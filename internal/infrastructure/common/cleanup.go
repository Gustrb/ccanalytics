package common

import "errors"

type CleanupFunction func() error

func JoinCleanup(cleanups []CleanupFunction) CleanupFunction {
	return func() error {
		var e error

		for _, cleanup := range cleanups {
			if err := cleanup(); err != nil {
				e = errors.Join(e, err)
			}
		}

		return e
	}
}
