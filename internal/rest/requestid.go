package rest

import (
	"context"
	"crypto/rand"
	"math/big"

	"github.com/Gustrb/ccanalytics/internal/rest/contextkey"
	"github.com/labstack/echo/v5"
)

var (
	genRandSet = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
)

func WithRequestID(next echo.HandlerFunc) echo.HandlerFunc {
	const headerKey = "X-Request-ID"

	genRand := func() (string, error) {
		r := make([]rune, 16)
		for i := range r {
			randomIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(genRandSet))))
			if err != nil {
				return "", err
			}

			r[i] = genRandSet[randomIdx.Int64()]
		}

		return string(r), nil
	}

	return func(c *echo.Context) error {
		requestID := c.Request().Header.Get(headerKey)
		if requestID == "" {
			r, err := genRand()
			if err != nil {
				return err
			}

			requestID = r
		}

		setContext(c, func(ctx context.Context) context.Context {
			return context.WithValue(ctx, contextkey.RequestIDKey, requestID)
		})

		c.Response().Header().Set(headerKey, requestID)

		return next(c)
	}
}
