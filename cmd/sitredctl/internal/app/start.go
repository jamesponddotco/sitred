package app

import (
	"fmt"
	"log/slog"
	"os"

	"git.sr.ht/~jamesponddotco/sitred/internal/config"
	"git.sr.ht/~jamesponddotco/sitred/internal/server"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/urfave/cli/v2"
)

// ErrServerRunning is returned when the server is already running.
const ErrServerRunning xerrors.Error = "server is already running"

// StartAction is the action for the start command.
func StartAction(ctx *cli.Context) error {
	cfg, err := config.Parse(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	srv, err := server.New(cfg, logger)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err = os.Stat(cfg.Server.PID); !os.IsNotExist(err) {
		return ErrServerRunning
	}

	pid := os.Getpid()

	pidFile, err := os.Create(cfg.Server.PID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer pidFile.Close()

	_, err = fmt.Fprintf(pidFile, "%d\n", pid)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	logger.LogAttrs(
		ctx.Context,
		slog.LevelInfo,
		"starting server",
		slog.String("pid", cfg.Server.PID),
		slog.String("address", cfg.Server.Address),
	)

	if err := srv.Start(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
