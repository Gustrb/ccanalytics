package binsign

import (
	"github.com/labstack/echo/v5"
)

func Urls(e *echo.Echo) {
	e.POST("/binsign/sign", SignHandler)
}
