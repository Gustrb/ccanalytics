package binsign

import (
	"context"
	"errors"

	"github.com/Gustrb/ccanalytics/internal/infrastructure/database"
)

const (
	insertSignedBinaryQuery = "insert into signed_binaries (hash, created_at, updated_at) values (?, ?, ?);"
)

var (
	ErrDuplicateHash = errors.New("a signed binary with the same hash already exists")
)

func Create(ctx context.Context, sb *SignedBinary) (*SignedBinary, error) {
	if err := database.InsertContext(ctx, insertSignedBinaryQuery, sb); err != nil {
		if database.IsDuplicateEntryError(err) {
			return nil, ErrDuplicateHash
		}

		return nil, err
	}

	return sb, nil
}
