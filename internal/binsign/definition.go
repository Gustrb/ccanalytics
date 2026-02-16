package binsign

import "time"

type SignedBinary struct {
	ID        int    `sql:"id"`
	Hash      string `sql:"hash"`
	CreatedAt int64  `sql:"created_at"`
	UpdatedAt int64  `sql:"updated_at"`
}

func (sb *SignedBinary) GetID() int {
	return sb.ID
}

func (sb *SignedBinary) SetID(id int) {
	sb.ID = id
}

type SignedBinaryOptions func(*SignedBinary)

func WithHash(hash string) SignedBinaryOptions {
	return func(sb *SignedBinary) {
		sb.Hash = hash
	}
}

func NewSignedBinary(opts ...SignedBinaryOptions) *SignedBinary {
	m := &SignedBinary{}

	for _, opt := range opts {
		opt(m)
	}

	m.CreatedAt = time.Now().UnixNano()
	m.UpdatedAt = time.Now().UnixNano()

	return m
}
