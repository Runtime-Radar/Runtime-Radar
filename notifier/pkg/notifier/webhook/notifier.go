package webhook

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/runtime-radar/runtime-radar/lib/security"
	"github.com/runtime-radar/runtime-radar/notifier/api"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/model"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/template"
)

const (
	timeout                 = time.Second * 5
	userAgentHeaderTemplate = "CS/%s"
	minTLSVersion           = tls.VersionTLS12
)

type Notifier struct {
	Config *model.Webhook
}

func newClient(conf *model.Webhook) (*http.Client, error) {
	var cas []string
	if conf.CA != "" {
		cas = []string{conf.CA}
	}

	cp, err := security.LoadCABundle(cas...)
	if err != nil {
		return nil, fmt.Errorf("can't load cert bundle: %w", err)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		MinVersion:         minTLSVersion,
		RootCAs:            cp,
		InsecureSkipVerify: conf.Insecure,
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}, nil
}

// Test tries to establish tcp connection with configured host.
// Optionally, tls is used in case when https scheme is presented in url.
func (n *Notifier) Test(ctx context.Context) error {
	u, err := url.Parse(n.Config.URL)
	if err != nil {
		return fmt.Errorf("can't parse url: %w", err)
	}

	host := u.Host
	nd := &net.Dialer{}

	dialCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if u.Scheme == "https" {
		if u.Port() == "" {
			host += ":443"
		}

		var cas []string
		if n.Config.CA != "" {
			cas = []string{n.Config.CA}
		}

		cp, err := security.LoadSystemCABundle(cas...)
		if err != nil {
			return fmt.Errorf("can't load cert bundle: %w", err)
		}

		td := tls.Dialer{
			NetDialer: nd,
			Config: &tls.Config{
				MinVersion:         minTLSVersion,
				RootCAs:            cp,
				InsecureSkipVerify: n.Config.Insecure,
			},
		}

		c, err := td.DialContext(dialCtx, "tcp", host)
		if err != nil {
			return fmt.Errorf("can't dial tls: %w", err)
		}
		defer c.Close()

		return nil
	}

	if u.Port() == "" {
		host += ":80"
	}

	c, err := nd.DialContext(dialCtx, "tcp", host)
	if err != nil {
		return fmt.Errorf("can't dial: %w", err)
	}
	defer c.Close()

	return nil
}

func (n *Notifier) Notify(ctx context.Context, notification *model.Notification, event any) error {
	c, err := newClient(n.Config)
	if err != nil {
		return fmt.Errorf("can't create client: %w", err)
	}

	body, err := renderRequestBody(notification, event)
	if err != nil {
		return fmt.Errorf("can't render body: %w", err)
	}

	u, err := buildURL(n.Config.URL, notification.WebhookConfig.Path)
	if err != nil {
		return fmt.Errorf("can't build url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}

	// currently json is only supported request body's format
	req.Header.Set("Content-Type", "application/json")

	if m := n.Config.Meta; m.CSVersion != "" {
		req.Header.Set("User-Agent", fmt.Sprintf(userAgentHeaderTemplate, m.CSVersion))
	}

	if wc := notification.WebhookConfig; wc != nil {
		for h, v := range wc.Headers {
			req.Header.Set(h, v)
		}
	}

	if n.Config.Login != "" {
		req.SetBasicAuth(n.Config.Login, n.Config.Password)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("can't send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("invalid status code returned: %d", resp.StatusCode)
	}

	return nil
}

func buildURL(base, path string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("can't parse base url: %w", err)
	}

	// Prevent query from being encoded as a part of path.
	path, query, _ := strings.Cut(path, "?")

	u = u.JoinPath(path)

	if query != "" {
		q, err := url.ParseQuery(query)
		if err != nil {
			return "", fmt.Errorf("can't parse query: %w", err)
		}

		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

// renderRequestBody parses and executes template if n has one. If it doesn't, event marshaled to json is returned.
func renderRequestBody(n *model.Notification, event any) (io.Reader, error) {
	var unwrappedEvent any

	switch ev := event.(type) {
	case *api.Message_RuntimeEvent:
		unwrappedEvent = ev.RuntimeEvent
	default:
		return nil, fmt.Errorf("invalid event type given: %T", ev)
	}

	tplData := map[string]any{
		"event":            unwrappedEvent,
		"notificationName": n.Name,
		"centralCSURL":     n.CentralCSURL,
		"csClusterID":      n.CSClusterID,
		"csClusterName":    n.CSClusterName,
		"ownCSURL":         n.OwnCSURL,
	}

	tpl, ok := template.DefaultTexts[n.EventType]
	if !ok {
		return nil, fmt.Errorf("no default template for '%s' event type", n.EventType)
	}

	if n.Template != "" {
		var err error

		tpl, err = template.NewText("", n.Template)
		if err != nil {
			return nil, fmt.Errorf("can't parse template: %w", err)
		}
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, tplData); err != nil {
		return nil, fmt.Errorf("can't execute template: %w", err)
	}

	var compactBuf bytes.Buffer
	if err := json.Compact(&compactBuf, buf.Bytes()); err != nil {
		return nil, fmt.Errorf("can't compact json: %w", err)
	}

	return &compactBuf, nil
}
