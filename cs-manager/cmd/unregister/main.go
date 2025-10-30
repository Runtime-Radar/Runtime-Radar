package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/cluster-manager/api"
	"github.com/runtime-radar/runtime-radar/cs-manager/pkg/client"
	"github.com/runtime-radar/runtime-radar/cs-manager/pkg/config"
	"github.com/runtime-radar/runtime-radar/lib/security"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
)

const (
	// TLS cert file name.
	certFile = "cert.pem"
	// TLS key file name.
	keyFile = "key.pem"
	// CA cert file name.
	caFile = "ca.pem"
)

func main() {
	cfg := config.New()

	centralCSURL, err := parseCSURL(cfg.CentralCSURL)
	if err != nil {
		log.Fatal().Msgf("### Failed to parse central CS URL: %v", err)
	}

	token, err := uuid.Parse(cfg.RegistrationToken)
	if err != nil {
		log.Fatal().Msgf("### Failed to parse registration token of current CS: %v", err)
	}

	var tokenKey []byte
	if cfg.Auth {
		_, tokenKey, err = jwt.NewKeyVerifier(cfg.TokenKey)
		if err != nil {
			log.Fatal().Msgf("### Failed to instantiate key verifier: %v", err)
		}
	}

	var tlsConfig *tls.Config
	if cfg.TLS {
		tlsConfig, err = security.LoadTLS(caFile, certFile, keyFile)
		if err != nil {
			log.Fatal().Msgf("### Failed to load TLS config: %v", err)
		}
		tlsConfig.InsecureSkipVerify = !cfg.CentralCSTLSCheckCert
	}

	clusterController, closeCC, err := client.NewClusterController(centralCSURL.Host, tlsConfig, tokenKey)
	if err != nil {
		log.Fatal().Msgf("### Failed to connect to central Cluster Manager: %v", err)
	}
	defer closeCC()

	_, err = clusterController.Unregister(context.Background(), &api.UnregisterClusterReq{
		Token: token.String(),
	})
	if err != nil {
		log.Fatal().Msgf("### Failed to deregister cluster: %v", err)
	}
}

func parseCSURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return &url.URL{}, nil
	}

	if !strings.Contains(rawURL, "://") {
		return nil, fmt.Errorf("wrong format, url should contain scheme: %s", rawURL)
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("can't parse url: %w", err)
	}

	// Explicitly set default port to avoid gRPC connection error when dialing with
	// a host ip address that doesn't include a port number
	// `transport: Error while dialing: dial tcp: address X.X.X.X: missing port in address`
	if parsedURL.Port() == "" {
		switch parsedURL.Scheme {
		case "http":
			parsedURL.Host = net.JoinHostPort(parsedURL.Hostname(), "80")
		case "https":
			parsedURL.Host = net.JoinHostPort(parsedURL.Hostname(), "443")
		default:
			return nil, fmt.Errorf("unsupported scheme: %s", rawURL)
		}
	}

	return parsedURL, nil
}
