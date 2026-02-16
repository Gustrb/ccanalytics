package binsign

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Gustrb/ccanalytics/internal/infrastructure/database"
)

const (
	getSignedBinaryByHashQuery = "select * from signed_binaries where hash = ?;"
)

var (
	ErrTooManyBinariesWithSameHash = fmt.Errorf("too many binaries with the same hash")
)

func GetSignedBinaryByHash(ctx context.Context, hash string) (*SignedBinary, error) {
	signedBinaries, err := database.SelectContext[SignedBinary](ctx, getSignedBinaryByHashQuery, hash)
	if err != nil {
		return nil, err
	}

	if len(signedBinaries) == 0 {
		slog.InfoContext(ctx, "No signed binary found with the given hash", "hash", hash)
		return nil, nil
	}

	if len(signedBinaries) > 1 {
		slog.WarnContext(ctx, "Multiple signed binaries found with the same hash, this should never happen", "hash", hash, "count", len(signedBinaries))
		return nil, fmt.Errorf("too many binaries with the same hash: %s, %w", hash, ErrTooManyBinariesWithSameHash)
	}

	return signedBinaries[0], nil
}
