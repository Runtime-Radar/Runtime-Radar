package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/runtime-radar/runtime-radar/lib/security"
	"github.com/runtime-radar/runtime-radar/notifier/api"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/model"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/template"
	"github.com/wneessen/go-mail"
)

const userAgentHeaderTemplate = "CS/%s"

const (
	dialTimeout           = 5 * time.Second
	testConnectionTimeout = 10 * time.Second
	sendConnectionTimeout = 5 * time.Minute
)

type Notifier struct {
	Config *model.Email
}

func newClient(e *model.Email, connTimeout time.Duration) (*mail.Client, error) {
	// Check if server has a port specified and split it if so
	server := e.Server
	port := mail.DefaultPort
	if strings.Contains(server, ":") {
		s, p, err := net.SplitHostPort(server)
		if err != nil {
			return nil, err
		}
		server = s
		if port, err = strconv.Atoi(p); err != nil {
			return nil, err
		}
	}

	var tlsConfig *tls.Config
	if e.UseTLS || e.UseStartTLS {
		var (
			cp  *x509.CertPool
			err error
			cas []string
		)

		if e.CA != "" {
			cas = []string{e.CA}
		}

		if cp, err = security.LoadSystemCABundle(cas...); err != nil {
			return nil, fmt.Errorf("can't load cert bundle: %w", err)
		}

		tlsConfig = &tls.Config{
			ServerName:         server,
			MinVersion:         mail.DefaultTLSMinVersion,
			RootCAs:            cp,
			InsecureSkipVerify: e.Insecure,
		}
	}

	// Due to there is only connection timeout and there isn't read/write
	// timeout in wneessen/go-mail library we need to process timeout
	// manually to avoid long queries.
	c, err := mail.NewClient(
		server,
		mail.WithPort(port),
		mail.WithDialContextFunc(func(ctx context.Context, network, address string) (net.Conn, error) {
			netDialer := net.Dialer{Timeout: dialTimeout}

			var conn net.Conn
			var err error
			if e.UseTLS {
				tlsDialer := tls.Dialer{
					NetDialer: &netDialer,
					Config:    tlsConfig,
				}
				conn, err = tlsDialer.DialContext(ctx, network, address)
			} else {
				conn, err = netDialer.DialContext(ctx, network, address)
			}
			if err != nil {
				return nil, fmt.Errorf("can't create dialer: %w", err)
			}

			if err = conn.SetDeadline(time.Now().Add(connTimeout)); err != nil {
				return nil, fmt.Errorf("can't set deadline: %w", err)
			}
			return conn, nil
		}))
	if err != nil {
		return nil, fmt.Errorf("can't create mail client: %w", err)
	}

	// Set username, password and auth type if auth is enabled
	if e.AuthType != model.AuthTypeNone {
		c.SetUsername(e.Username)
		c.SetPassword(e.Password)
		c.SetSMTPAuth(smtpAuthTypeFromEnum(e.AuthType))
	}

	// UseTLS initiates a TLS connection to the server without STARTTLS
	c.SetSSL(e.UseTLS)

	c.SetTLSPolicy(mail.NoTLS)
	if e.UseStartTLS {
		c.SetTLSPolicy(mail.TLSMandatory)
	}

	if tlsConfig != nil {
		if err := c.SetTLSConfig(tlsConfig); err != nil {
			return nil, fmt.Errorf("can't set tls config: %w", err)
		}
	}

	return c, nil
}

// Test tests connection and auth credentials for SMTP-server.
func (n *Notifier) Test(ctx context.Context) error {
	c, err := newClient(n.Config, testConnectionTimeout)
	if err != nil {
		return err
	}
	if err := c.DialWithContext(ctx); err != nil {
		return err
	}
	// c.Close() call must be after succeeded Dial*() method call
	// because c.Close() uses internal *smtp.Client variable
	// that is set up when using the Dial*() methods
	return c.Close()
}

func (n *Notifier) Notify(ctx context.Context, notification *model.Notification, event any) error {
	if len(notification.Recipients) == 0 {
		return fmt.Errorf("no recipients")
	}

	var (
		subject, text string
		err           error
	)

	switch ev := event.(type) {
	case *api.Message_RuntimeEvent:
		subject, text, err = messageFromRuntimeMonitorEvent(notification, ev.RuntimeEvent)
	default:
		return fmt.Errorf("invalid event type given: %w", err)
	}

	if err != nil {
		return fmt.Errorf("can't build message from event: %w", err)
	}

	m := mail.NewMsg()
	if err := m.From(n.Config.From); err != nil {
		return fmt.Errorf("can't set From: %w", err)
	}

	if err := m.To(notification.Recipients...); err != nil {
		return fmt.Errorf("can't set To: %w", err)
	}

	if meta := n.Config.Meta; meta.CSVersion != "" {
		m.SetUserAgent(fmt.Sprintf(userAgentHeaderTemplate, meta.CSVersion))
	}
	m.SetBodyString(mail.TypeTextHTML, text)
	m.Subject(subject)

	c, err := newClient(n.Config, sendConnectionTimeout)
	if err != nil {
		return fmt.Errorf("can't create client: %w", err)
	}

	return c.DialAndSendWithContext(ctx, m)
}

func messageFromRuntimeMonitorEvent(n *model.Notification, event *api.RuntimeEvent) (subject, text string, err error) {
	if conf := n.EmailConfig; conf != nil && conf.SubjectTemplate != "" {
		subject, err = renderSubject(n, event)
		if err != nil {
			return "", "", fmt.Errorf("can't render subject: %w", err)
		}
	} else {
		subject = "CS: threats detected in the runtime event."
	}

	text, err = renderText(n, event)
	if err != nil {
		return "", "", err
	}

	return subject, text, nil
}

// renderSubjects parses and executes subject template. It assumes that n.EmailConfig is not nil.
func renderSubject(n *model.Notification, event any) (string, error) {
	t, err := template.NewText("", n.EmailConfig.SubjectTemplate)
	if err != nil {
		return "", fmt.Errorf("can't parse template: %w", err)
	}

	templateData := map[string]any{
		"event":            event,
		"notificationName": n.Name,
		"centralCSURL":     n.CentralCSURL,
		"csClusterID":      n.CSClusterID,
		"csClusterName":    n.CSClusterName,
		"ownCSURL":         n.OwnCSURL,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("can't execute template: %w", err)
	}

	return buf.String(), nil
}

func renderText(n *model.Notification, event any) (text string, err error) {
	tpl, ok := template.DefaultHTMLs[n.EventType]
	if !ok {
		return "", fmt.Errorf("no default template for %s event type", n.EventType)
	}

	if n.Template != "" {
		tpl, err = template.NewHTML("", n.Template)
		if err != nil {
			return "", fmt.Errorf("can't parse template: %w", err)
		}
	}

	templateData := map[string]any{
		"event":            event,
		"notificationName": n.Name,
		"centralCSURL":     n.CentralCSURL,
		"csClusterID":      n.CSClusterID,
		"csClusterName":    n.CSClusterName,
		"ownCSURL":         n.OwnCSURL,
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("can't execute template: %w", err)
	}

	return buf.String(), nil
}

func smtpAuthTypeFromEnum(a model.EmailAuthType) mail.SMTPAuthType {
	switch a {
	case model.AuthTypePlain:
		return mail.SMTPAuthPlain
	case model.AuthTypeLogin:
		return mail.SMTPAuthLogin
	case model.AuthTypeCramMD5:
		return mail.SMTPAuthCramMD5
	default: // normally should not happen
		panic(fmt.Sprintf("invalid AuthType value: %d", a))
	}
}
