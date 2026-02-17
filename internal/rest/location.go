package rest

import (
	"context"
	"strconv"
	"strings"

	"github.com/Gustrb/ccanalytics/internal/rest/contextkey"
	"github.com/Gustrb/ccanalytics/internal/rest/location"
	"github.com/labstack/echo/v5"
)

func WithLocation(next echo.HandlerFunc) echo.HandlerFunc {
	parseBool := func(value string) *bool {
		value = strings.TrimSpace(value)

		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil
		}

		return &boolVal
	}

	return func(c *echo.Context) error {
		header := c.Request().Header

		loc := location.Location{
			City:        strings.TrimSpace(header.Get("Cloudfront-Viewer-City")),
			CountryCode: strings.TrimSpace(header.Get("Cloudfront-Viewer-Country")),
			CountryName: strings.TrimSpace(header.Get("Cloudfront-Viewer-Country-Name")),
			TimeZone:    strings.TrimSpace(header.Get("Cloudfront-Viewer-Time-Zone")),
			IsDesktop:   parseBool(header.Get("Cloudfront-Is-Desktop-Viewer")),
			IsMobile:    parseBool(header.Get("Cloudfront-Is-Mobile-Viewer")),
			IsTablet:    parseBool(header.Get("Cloudfront-Is-Tablet-Viewer")),
		}

		setContext(c, func(ctx context.Context) context.Context {
			return context.WithValue(ctx, contextkey.LocationKey, &loc)
		})

		return next(c)
	}
}
