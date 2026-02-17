package rest

import (
	"context"
	"slices"

	"github.com/Gustrb/ccanalytics/internal/infrastructure/database"
	"github.com/labstack/echo/v5"
)

var whitelistedRoutes = []string{
	"/",
	"/health",
}

func WithTransaction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if slices.Contains(whitelistedRoutes, c.Path()) {
			return next(c)
		}

		err := database.WithinTransaction(c.Request().Context(), func(ctx context.Context) error {
			// exec the next handler with the transaction context
			if err := next(c); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return echo.NewHTTPError(500, "An error occurred while processing the request").Wrap(err)
		}

		return nil
	}
}
