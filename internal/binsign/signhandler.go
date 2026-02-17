package binsign

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func SignHandler(c *echo.Context) error {
	// TODO: check idempotency key
	fheader, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "file is required")
	}

	fileHandle, err := fheader.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open file")
	}
	defer fileHandle.Close()

	ctx := c.Request().Context()

	if err := SignFile(ctx, fileHandle); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to sign file")
	}

	return nil
}
