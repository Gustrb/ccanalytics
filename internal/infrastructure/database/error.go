package database

import "strings"

func IsNoTableFoundError(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "no such table") {
		return true
	}

	return false
}
