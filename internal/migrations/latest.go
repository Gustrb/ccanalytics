package migrations

import (
	"context"

	"github.com/Gustrb/ccanalytics/internal/infrastructure/database"
)

const (
	getLatestAppliedMigrationsQuery = "select * from migrations order by timestamp;"
	getLastAppliedMigrationQuery    = "select * from migrations order by timestamp desc limit 1;"
)

func GetLatestAppliedMigrations(ctx context.Context) ([]*Migration, error) {
	migrations, err := database.SelectContext[Migration](ctx, getLatestAppliedMigrationsQuery)

	// If the err is no table found, we can ignore it and return an empty slice of migrations
	if database.IsNoTableFoundError(err) {
		return []*Migration{}, nil
	}
	if err != nil {
		return nil, err
	}

	return migrations, err
}

func GetLastMigration(ctx context.Context) (*Migration, error) {
	migrations, err := database.SelectContext[Migration](ctx, getLastAppliedMigrationQuery)

	// If the err is no table found, we can ignore it and return an empty slice of migrations
	if database.IsNoTableFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if len(migrations) == 0 {
		return nil, nil
	}

	return migrations[0], nil
}
