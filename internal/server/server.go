// Package server provides a simple HTTP server for the service.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sr.ht/~jamesponddotco/sitred/internal/config"
	"git.sr.ht/~jamesponddotco/sitred/internal/endpoint"
	"git.sr.ht/~jamesponddotco/sitred/internal/fetch"
	"git.sr.ht/~jamesponddotco/sitred/internal/server/handler"
	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xtls"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp/xmiddleware"
)

const (
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	DefaultIdleTimeout  = 60 * time.Second
)

// Server represents a Privytar server.
type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

// New creates a new HTTP server.
func New(cfg *config.Config, logger *slog.Logger) (*Server, error) {
	cert, err := tls.LoadX509KeyPair(cfg.Server.TLS.Certificate, cfg.Server.TLS.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	var tlsConfig *tls.Config

	if cfg.Server.TLS.Version == "1.3" {
		tlsConfig = xtls.ModernServerConfig()
	}

	if cfg.Server.TLS.Version == "1.2" {
		tlsConfig = xtls.IntermediateServerConfig()
	}

	tlsConfig.Certificates = []tls.Certificate{cert}

	middlewares := []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler { return xmiddleware.PanicRecovery(logger, h) },
		func(h http.Handler) http.Handler { return xmiddleware.UserAgent(logger, h) },
		func(h http.Handler) http.Handler {
			return xmiddleware.AcceptRequests(
				[]string{
					http.MethodGet,
					http.MethodHead,
					http.MethodOptions,
				},
				logger,
				h,
			)
		},
	}

	if cfg.Server.LogRequests {
		accessLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		middlewares = append(middlewares, func(h http.Handler) http.Handler {
			return xmiddleware.AccessLog(accessLogger, h)
		})
	}

	var (
		fetchInstance = fetch.New(cfg.Service.Name, cfg.Service.Contact)
		rootHandler   = handler.NewRootHandler(fetchInstance, logger, cfg.Sitemap.URL)
	)

	mux := http.NewServeMux()
	mux.Handle(endpoint.Root, xmiddleware.Chain(rootHandler, middlewares...))

	httpServer := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

// Start starts the Privytar server.
func (s *Server) Start() error {
	var (
		sigint            = make(chan os.Signal, 1)
		shutdownCompleted = make(chan struct{})
	)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.LogAttrs(
				ctx,
				slog.LevelError,
				"failed to shutdown server",
				slog.String("error", err.Error()),
			)
		}

		close(shutdownCompleted)
	}()

	if err := s.httpServer.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	<-shutdownCompleted

	return nil
}

// Stop gracefully shuts down the Privytar server.
func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
