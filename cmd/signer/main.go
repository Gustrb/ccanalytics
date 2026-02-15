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
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file_path",
				Required: true,
				Usage:    "a valid file path is required in order to be signed",
			},
			&cli.Uint16Flag{
				Name:  "timeout",
				Usage: "the duration in seconds to wait before timing out the signing process",
				Value: 1, // default timeout of 1 second
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			timeout := c.Uint16("timeout")
			slog.InfoContext(ctx, "Starting signer command")

			ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()

			select {
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					slog.WarnContext(ctx, "Signing process timed out")
				} else {
					slog.InfoContext(ctx, "Signing process completed successfully")
				}

			default:
				if err := binsign.SignFileAt(ctx, c.String("file_path")); err != nil {
					slog.ErrorContext(ctx, "Failed to sign file", "error", err)
					return err
				}
			}

			return nil
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cleanup, err := cmdutils.SetupBinary(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to set up signer command", "error", err)
		return
	}
	defer func() {
		if err := cleanup(); err != nil {
			slog.ErrorContext(ctx, "Failed to clean up resources", "error", err)
		}
	}()

	if err := cmd.Run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "Failed to run signer command", "error", err)
	}
}
