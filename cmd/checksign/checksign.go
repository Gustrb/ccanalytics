package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gustrb/ccanalytics/internal/binsign"
	"github.com/Gustrb/ccanalytics/internal/cmdutils"
	"github.com/Gustrb/ccanalytics/internal/infrastructure/database/migrator"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file_path",
				Required: true,
				Usage:    "a valid file path is required in order to be checked",
			},
			&cli.Uint16Flag{
				Name:  "timeout",
				Usage: "the duration in seconds to wait before timing out the checksigning process",
				Value: 1, // default timeout of 1 second
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			timeout := c.Uint16("timeout")
			slog.InfoContext(ctx, "Starting checksign command")

			ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()

			select {
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					slog.WarnContext(ctx, "Checksigning process timed out")
				} else {
					slog.InfoContext(ctx, "Checksigning process completed successfully")
				}

			default:
				signedBinary, err := binsign.CheckIfFileIsSigned(ctx, c.String("file_path"))
				if err != nil {
					slog.ErrorContext(ctx, "Failed to check if file is signed", "error", err)
					return err
				}

				if signedBinary == nil {
					slog.InfoContext(ctx, "File is not signed", "file_path", c.String("file_path"))
				} else {
					signedAt := time.Unix(0, signedBinary.CreatedAt).Format(time.RFC3339)
					slog.InfoContext(ctx, "File is signed", "file_path", c.String("file_path"), "signed_at", signedAt)
				}
			}

			return nil
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cleanup, err := cmdutils.SetupBinary(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to set up checksign command", "error", err)
		return
	}
	defer func() {
		if err := cleanup(); err != nil {
			slog.ErrorContext(ctx, "Failed to clean up resources", "error", err)
		}
	}()

	areWeAtTheLatestMigration, err := migrator.AreWeAtTheLatestMigration(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check for latest migration", "error", err)
		return
	}

	if !areWeAtTheLatestMigration {
		slog.WarnContext(ctx, "Database is not at the latest migration. Please run the migrator command to apply all pending migrations before running the checksign command.")
		return
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "Failed to run checksign command", "error", err)
	}
}
