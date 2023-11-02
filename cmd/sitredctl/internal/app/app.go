// Package app is the main package for the control application.
package app

import (
	"context"
	"log/slog"
	"os"

	"git.sr.ht/~jamesponddotco/sitred"
	"git.sr.ht/~jamesponddotco/sitred/internal/config"
	"github.com/urfave/cli/v2"
)

// logger is the default structured logger for the application.
var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil)) //nolint:gochecknoglobals // I have no idea how to pass the logger to a cli.Action otherwise.

// Run is the entry point for the application.
func Run(args []string) int {
	app := cli.NewApp()
	app.Name = sitred.Name
	app.Version = sitred.Version
	app.Usage = sitred.Description
	app.HideHelpCommand = true

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "server-pid",
			Usage: "path to pid file",
			Value: config.DefaultPID,
			EnvVars: []string{
				sitred.EnvPrefix + "_SERVER_PID",
			},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "start",
			Usage:  "start the server",
			Action: StartAction,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "tls-certificate",
					Usage: "path to TLS certificate",
					EnvVars: []string{
						sitred.EnvPrefix + "_TLS_CERTIFICATE",
					},
					Required: true,
				},
				&cli.StringFlag{
					Name:  "tls-key",
					Usage: "path to TLS key",
					EnvVars: []string{
						sitred.EnvPrefix + "_TLS_KEY",
					},
					Required: true,
				},
				&cli.StringFlag{
					Name:  "tls-version",
					Usage: "minimum TLS version supported by the server",
					Value: config.DefaultMinTLSVersion,
					EnvVars: []string{
						sitred.EnvPrefix + "_TLS_VERSION",
					},
				},
				&cli.StringFlag{
					Name:  "server-address",
					Usage: "address to bind to",
					Value: config.DefaultAddress,
					EnvVars: []string{
						sitred.EnvPrefix + "_SERVER_ADDRESS",
					},
				},
				&cli.DurationFlag{
					Name:  "server-cache-ttl",
					Usage: "time-to-live for cache entries",
					Value: config.DefaultCacheTTL,
					EnvVars: []string{
						sitred.EnvPrefix + "_CACHE_TTL",
					},
				},
				&cli.BoolFlag{
					Name:  "server-access-log",
					Usage: "whether to enable access logs",
					Value: false,
					EnvVars: []string{
						sitred.EnvPrefix + "_ACCESS_LOG",
					},
				},
				&cli.StringFlag{
					Name:  "service-name",
					Usage: "name of the service",
					Value: sitred.Name,
					EnvVars: []string{
						sitred.EnvPrefix + "_SERVICE_NAME",
					},
				},
				&cli.StringFlag{
					Name:  "service-contact",
					Usage: "contact information for the service",
					Value: sitred.URL,
					EnvVars: []string{
						sitred.EnvPrefix + "_SERVICE_CONTACT",
					},
				},
				&cli.StringFlag{
					Name:  "sitemap-url",
					Usage: "url of the sitemap to use when choosing random URLs",
					EnvVars: []string{
						sitred.EnvPrefix + "_SITEMAP_URL",
					},
					Required: true,
				},
			},
		},
		{
			Name:   "stop",
			Usage:  "stop the server",
			Action: StopAction,
		},
	}

	if err := app.Run(args); err != nil {
		logger.LogAttrs(
			context.Background(),
			slog.LevelError,
			"failed to initialize control application",
			slog.String("error", err.Error()),
		)

		return 1
	}

	return 0
}
