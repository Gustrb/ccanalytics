package binsign

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

func CheckSignHandler(c *echo.Context) error {
	ctx := c.Request().Context()

	fheader, err := c.FormFile("file")
	if err != nil {
		slog.ErrorContext(ctx, "failed to get file from form data", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "file is required")
	}

	fileHandle, err := fheader.Open()
	if err != nil {
		slog.ErrorContext(ctx, "failed to open file", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open file")
	}
	defer fileHandle.Close()

	signedFile, err := CheckIfReaderIsSigned(ctx, fileHandle)
	if err != nil {
		slog.ErrorContext(ctx, "failed to check if file is signed", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check if file is signed")
	}

	if signedFile == nil {
		// 404
		slog.InfoContext(ctx, "file is not signed", "file_name", fheader.Filename)
		return echo.NewHTTPError(http.StatusNotFound, "file is not signed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"file_name": fheader.Filename,
		"signed_at": signedFile.CreatedAt,
	})
}
