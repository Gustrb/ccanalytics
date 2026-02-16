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

func IsDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return true
	}

	return false
}
