package database

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"fmt"

	"github.com/Gustrb/ccanalytics/internal/infrastructure/common"
)

var (
	db *sql.DB
)

var (
	ErrDBAlreadyConnected = fmt.Errorf("database connection already established")
)

func Connect(ctx context.Context, dataSourceName string) (common.CleanupFunction, error) {
	if db != nil {
		return nil, ErrDBAlreadyConnected
	}

	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	cleanup := func() error {
		return db.Close()
	}

	if err := db.PingContext(ctx); err != nil {
		if cerr := cleanup(); cerr != nil {
			return nil, fmt.Errorf("failed to ping database: %w; also failed to clean up database connection: %v", err, cerr)
		}

		return nil, err
	}

	return cleanup, nil
}
