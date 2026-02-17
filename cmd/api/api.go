package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gustrb/ccanalytics/internal/cmdutils"
	"github.com/Gustrb/ccanalytics/internal/config"
	"github.com/Gustrb/ccanalytics/internal/handlers"
	"github.com/Gustrb/ccanalytics/internal/infrastructure/database"
	"github.com/Gustrb/ccanalytics/internal/infrastructure/database/migrator"
	"github.com/Gustrb/ccanalytics/internal/rest"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
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

	e := echo.New()
	recoverConfig := middleware.DefaultRecoverConfig

	if config.Environments.EnviromnmentName == config.EnvironmentProduction {
		recoverConfig.DisablePrintStack = true
		recoverConfig.DisableStackAll = true
	}

	e.Use(middleware.RecoverWithConfig(recoverConfig))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS("*"))
	e.Use(rest.WithTransaction)
	e.Use(rest.WithRequestID)
	e.Use(rest.WithLocation)
	e.Use(rest.WithLogging)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{})
	})

	e.GET("/health", func(c *echo.Context) error {
		ctx, cancel := context.WithDeadline(c.Request().Context(), time.Now().Add(time.Millisecond*200))
		defer cancel()

		if err := database.PingContext(ctx); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database is not healthy").Wrap(err)
		}

		return c.JSON(http.StatusOK, map[string]string{})
	})

	handlers.Register(e)

	slog.InfoContext(ctx, "Starting web server", "addr", config.Rest.Addr)

	if err := e.Start(config.Rest.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.ErrorContext(ctx, "Failed to start web server", "error", err)
		return
	}
}
