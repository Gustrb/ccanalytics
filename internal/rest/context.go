package rest

import (
	"context"

	"github.com/labstack/echo/v5"
)

func setContext(c *echo.Context, f func(context.Context) context.Context) {
	rA := c.Request()
	rB := rA.WithContext(f(rA.Context()))
	c.SetRequest(rB)
}
