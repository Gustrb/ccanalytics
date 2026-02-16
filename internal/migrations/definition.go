package migrations

import "time"

type Migration struct {
	ID        int    `sql:"id"`
	Filename  string `sql:"filename"`
	Timestamp int64  `sql:"timestamp"`
	CreatedAt int64  `sql:"created_at"`
	UpdatedAt int64  `sql:"updated_at"`
}

func (m *Migration) GetID() int {
	return m.ID
}

func (m *Migration) SetID(id int) {
	m.ID = id
}

type MigrationOptions func(*Migration)

func WithFilename(filename string) MigrationOptions {
	return func(m *Migration) {
		m.Filename = filename
	}
}

func WithTimestamp(timestamp int64) MigrationOptions {
	return func(m *Migration) {
		m.Timestamp = timestamp
	}
}

func NewMigration(opts ...MigrationOptions) *Migration {
	m := &Migration{}

	for _, opt := range opts {
		opt(m)
	}

	m.CreatedAt = time.Now().UnixNano()
	m.UpdatedAt = time.Now().UnixNano()

	return m
}
