package handlers

import (
	"github.com/Gustrb/ccanalytics/internal/binsign"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo) {
	binsign.Urls(e)
}
