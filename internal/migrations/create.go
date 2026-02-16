package migrations

import (
	"context"

	"github.com/Gustrb/ccanalytics/internal/infrastructure/database"
)

const (
	insertMigrationQuery = "insert into migrations (filename, timestamp, created_at, updated_at) values (?, ?, ?, ?);"
)

func Create(ctx context.Context, migration *Migration) (*Migration, error) {
	if err := database.InsertContext(ctx, insertMigrationQuery, migration); err != nil {
		return nil, err
	}

	return migration, nil
}
