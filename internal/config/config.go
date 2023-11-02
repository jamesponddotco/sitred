// Package config implements the configuration logic for the service.
package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"git.sr.ht/~jamesponddotco/sitred"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/urfave/cli/v2"
)

const (
	// ErrInvalidConfig is returned when any Config field is invalid.
	ErrInvalidConfig xerrors.Error = "invalid configuration"

	// ErrMissingTLSCertificate is returned when the TLS certificate is missing.
	ErrMissingTLSCertificate xerrors.Error = "server's TLS certificate is missing"

	// ErrMissingTLSKey is returned when the TLS key is missing.
	ErrMissingTLSKey xerrors.Error = "server's TLS key is missing"

	// ErrInvalidTLSVersion is returned when the TLS version is invalid.
	ErrInvalidTLSVersion xerrors.Error = "server's TLS version is invalid; must be 1.2 or 1.3"

	// ErrMissingServerAddress is returned when the server address is missing.
	ErrMissingServerAddress xerrors.Error = "server address is missing"

	// ErrMissingServerPID is returned when the server PID is missing.
	ErrMissingServerPID xerrors.Error = "server PID is missing"

	// ErrMissingServiceName is returned when the service name is missing.
	ErrMissingServiceName xerrors.Error = "service name is missing"

	// ErrMissingServiceContact is returned when the service contact is missing.
	ErrMissingServiceContact xerrors.Error = "service contact is missing"

	// ErrMissingSitemapURL is returned when the sitemap URL is missing.
	ErrMissingSitemapURL xerrors.Error = "sitemap URL is missing"

	// ErrInvalidServerCacheTTL is returned when the server cache TTL is invalid.
	ErrInvalidServerCacheTTL xerrors.Error = "server cache TTL is invalid; must be a positive duration"

	// ErrInvalidSitemapURL is returned when the sitemap URL is invalid.
	ErrInvalidSitemapURL xerrors.Error = "sitemap URL is invalid; must be a valid URL and cannot be an index page"
)

const (
	// DefaultMinTLSVersion is the default minimum TLS version supported by the
	// server.
	DefaultMinTLSVersion string = "1.3"

	// DefaultAddress is the default address of the application.
	DefaultAddress string = ":1997"

	// DefaultPID is the default path to the PID file.
	DefaultPID string = "/var/run/" + sitred.Name + ".pid"

	// DefaultCacheTTL is the default time-to-live of the cache.
	DefaultCacheTTL time.Duration = 30 * time.Minute

	// DefaultServiceName is the default name of the service.
	DefaultServiceName string = sitred.Name
)

// TLS represents the TLS configuration.
type TLS struct {
	// Certificate is the path to the TLS certificate.
	Certificate string

	// Key is the path to the TLS key.
	Key string

	// Version is the TLS version to use.
	Version string
}

// Server represents the server configuration.
type Server struct {
	// TLS is the TLS configuration.
	TLS *TLS

	// Address is the address of the application.
	Address string

	// PID is the path to the PID file.
	PID string

	// CacheTTL is the TTL of the cache.
	CacheTTL time.Duration

	// LogRequests defines whether the application should log requests.
	LogRequests bool
}

// Service represents the service configuration.
type Service struct {
	// Name is the name of the service.
	Name string

	// Contact is the contact address for the service.
	Contact string
}

// Sitemap represents the sitemap configuration.
type Sitemap struct {
	// URL is the URL of the sitemap.
	URL string
}

// Config represents the application configuration.
type Config struct {
	// Service is the service configuration.
	Service *Service

	// Server is the server configuration.
	Server *Server

	// Sitemap is the sitemap configuration.
	Sitemap *Sitemap
}

// Parse parses a cli.Context and returns a Config from it or an error if the
// Config isn't valid.
func Parse(ctx *cli.Context) (*Config, error) {
	cfg := &Config{
		Service: &Service{
			Name:    ctx.String("service-name"),
			Contact: ctx.String("service-contact"),
		},
		Server: &Server{
			TLS: &TLS{
				Certificate: ctx.String("tls-certificate"),
				Key:         ctx.String("tls-key"),
				Version:     ctx.String("tls-version"),
			},
			Address:     ctx.String("server-address"),
			PID:         ctx.String("server-pid"),
			CacheTTL:    ctx.Duration("server-cache-ttl"),
			LogRequests: ctx.Bool("server-log-requests"),
		},
		Sitemap: &Sitemap{
			URL: ctx.String("sitemap-url"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	return cfg, nil
}

// Validate checks Config for errors.
func (cfg *Config) Validate() error {
	if cfg.Service.Name == "" {
		return ErrMissingServiceName
	}

	if cfg.Service.Contact == "" {
		return ErrMissingServiceContact
	}

	if cfg.Server.TLS.Certificate == "" {
		return ErrMissingTLSCertificate
	}

	if cfg.Server.TLS.Key == "" {
		return ErrMissingTLSKey
	}

	if cfg.Server.TLS.Version != "1.2" && cfg.Server.TLS.Version != "1.3" {
		return ErrInvalidTLSVersion
	}

	if cfg.Server.Address == "" {
		return ErrMissingServerAddress
	}

	if cfg.Server.PID == "" {
		return ErrMissingServerPID
	}

	if cfg.Server.CacheTTL == 0 {
		return ErrInvalidServerCacheTTL
	}

	if cfg.Sitemap.URL == "" {
		return ErrMissingSitemapURL
	}

	// Checks for the suffix used by Yoast SEO for WordPress. This is a naive
	// check and only works for that specific case, but it's good enough for
	// something this early.
	if strings.HasSuffix(cfg.Sitemap.URL, "sitemap_index.xml") {
		return ErrInvalidSitemapURL
	}

	if _, err := url.Parse(cfg.Sitemap.URL); err != nil {
		return ErrInvalidSitemapURL
	}

	return nil
}
