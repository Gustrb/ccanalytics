package rest

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
)

func WithLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		startTime := time.Now()

		ctx := c.Request().Context()

		err := next(c)
		if err != nil {
			c.Response().WriteHeader(500)
		}

		if errors.Is(err, context.Canceled) {
			return nil
		}

		// handler may have changed ctx, re-request
		ctx = c.Request().Context()

		attrs := []any{
			slog.Int64("response.elapsed", time.Since(startTime).Milliseconds()),
		}

		slog.InfoContext(ctx, "Returning response", attrs...)

		return nil
	}
}
