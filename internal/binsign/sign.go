package binsign

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// Note: crazy early optimization, maybe it is a silly one. But I wanted to use sync.Pools. Remove if it doesn't make a difference.
var bufferPool = sync.Pool{
	New: func() any {
		return make([]byte, 32) // SHA-256 produces a 32-byte hash
	},
}

func SignFileAt(ctx context.Context, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	hash, err := digest(f)
	if err != nil {
		return err
	}

	signedBinary := NewSignedBinary(
		WithHash(hash),
	)
	if _, err := Create(ctx, signedBinary); err != nil {
		return err
	}

	return nil
}

func CheckIfFileIsSigned(ctx context.Context, filePath string) (*SignedBinary, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hash, err := digest(f)
	if err != nil {
		return nil, err
	}

	signedBinary, err := GetSignedBinaryByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	if signedBinary == nil {
		slog.InfoContext(ctx, "No signed binary found with the given hash")
	} else {
		signedAt := time.Unix(signedBinary.CreatedAt, 0).Format(time.RFC3339)
		slog.InfoContext(ctx, "The binary was signed", "signed_at", signedAt)
	}

	return signedBinary, nil
}

func digest(reader io.Reader) (string, error) {
	hasher := sha256.New()

	if _, err := io.Copy(hasher, reader); err != nil {
		return "", err
	}

	buff, ok := bufferPool.Get().([]byte)
	if !ok {
		buff = make([]byte, 32)
	}

	hashBytes := hasher.Sum(buff)

	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
