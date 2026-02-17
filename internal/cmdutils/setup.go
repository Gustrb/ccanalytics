package cmdutils

import (
	"context"

	_ "github.com/Gustrb/ccanalytics/internal/config"
	"github.com/Gustrb/ccanalytics/internal/infrastructure/common"
	"github.com/Gustrb/ccanalytics/internal/infrastructure/database"
)

var (
	dbFilePath = "app.db"
)

func SetupBinary(ctx context.Context) (common.CleanupFunction, error) {
	cleanups := []common.CleanupFunction{}

	cleanup, err := database.Connect(ctx, dbFilePath)
	if err != nil {
		return nil, err
	}

	cleanups = append(cleanups, cleanup)

	return common.JoinCleanup(cleanups), nil
}
