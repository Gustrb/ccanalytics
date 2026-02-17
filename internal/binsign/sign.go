package binsign

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sync"
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

	return SignFile(ctx, f)
}

func SignFile(ctx context.Context, reader io.Reader) error {
	hash, err := digest(reader)
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

	return CheckIfReaderIsSigned(ctx, f)
}

func CheckIfReaderIsSigned(ctx context.Context, reader io.Reader) (*SignedBinary, error) {
	hash, err := digest(reader)
	if err != nil {
		return nil, err
	}

	signedBinary, err := GetSignedBinaryByHash(ctx, hash)
	if err != nil {
		return nil, err
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
