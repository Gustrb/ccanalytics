package migrator

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Gustrb/ccanalytics/internal/infrastructure/database"
	"github.com/Gustrb/ccanalytics/internal/migrations"
)

var (
	//go:embed migrations
	migrationsFS embed.FS
)

var (
	timestampReg, _ = regexp.Compile(`^\d+`)
)

type Migration struct {
	Up        string
	Down      string
	Filename  string
	Timestamp int
}

func MigrateUp(ctx context.Context) error {
	changes, err := parseMigrations()
	if err != nil {
		return err
	}

	latestMigrations, err := migrations.GetLatestAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	if len(latestMigrations) == len(changes) {
		slog.InfoContext(ctx, "All migrations have already been applied")
		return nil
	}

	if len(latestMigrations) > len(changes) {
		return fmt.Errorf("more migrations have been applied than exist in the codebase, this should never happen")
	}

	if len(latestMigrations) == 0 {
		slog.InfoContext(ctx, "No migrations have been applied yet, applying all migrations")
		if err := applyMigrations(ctx, changes); err != nil {
			return fmt.Errorf("applying migrations: %w", err)
		}
		return nil
	}

	// Now we need to find the last common migration and apply the rest...
	lastCommonIndex := -1
	for i, change := range changes {
		if i >= len(latestMigrations) {
			break
		}

		appliedMigration := latestMigrations[i]

		if change.Timestamp != int(appliedMigration.Timestamp) {
			break
		}

		lastCommonIndex = i
	}

	if lastCommonIndex == -1 {
		return fmt.Errorf("no common migration found, this should never happen")
	}

	if lastCommonIndex == len(changes)-1 {
		slog.InfoContext(ctx, "all migrations have already been applied")
	}

	if err := applyMigrations(ctx, changes[lastCommonIndex+1:]); err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}

func applyMigrations(ctx context.Context, migrationList []*Migration) error {
	err := database.WithinTransaction(ctx, func(ctx context.Context) error {
		for _, migration := range migrationList {
			slog.InfoContext(ctx, "applying migration", "filename", migration.Filename)

			resultSet, err := database.ExecContext(ctx, migration.Up)
			if err != nil {
				return fmt.Errorf("applying migration %s: %w", migration.Filename, err)
			}

			rowsAffected, err := resultSet.RowsAffected()
			if err != nil {
				return fmt.Errorf("getting rows affected for migration %s: %w", migration.Filename, err)
			}

			m := migrations.NewMigration(
				migrations.WithFilename(migration.Filename),
				migrations.WithTimestamp(int64(migration.Timestamp)),
			)
			if _, err := migrations.Create(ctx, m); err != nil {
				return fmt.Errorf("inserting migration record for %s: %w", migration.Filename, err)
			}

			slog.InfoContext(ctx, "applied migration", "filename", migration.Filename, "rowsAffected", rowsAffected)
		}
		return nil
	})
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully applied all migrations")

	return nil
}

func AreWeAtTheLatestMigration(ctx context.Context) (bool, error) {
	changes, err := parseMigrations()
	if err != nil {
		return false, err
	}

	lastMigration, err := migrations.GetLastMigration(ctx)
	if err != nil {
		return false, err
	}
	if lastMigration == nil {
		return len(changes) == 0, nil
	}

	return changes[len(changes)-1].Timestamp == int(lastMigration.Timestamp), nil
}

func parseMigrations() ([]*Migration, error) {
	filenames, err := fs.Glob(migrationsFS, "*/*.sql")
	if err != nil {
		return nil, fmt.Errorf("globbing migrations: %w", err)
	}

	changes := make([]*Migration, 0, len(filenames))

	for _, filename := range filenames {
		readFile := func() error {
			file, err := migrationsFS.Open(filename)
			if err != nil {
				return err
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			m, err := parseMigration(data, filename)
			if err != nil {
				return err
			}

			changes = append(changes, m)

			return nil
		}

		if err := readFile(); err != nil {
			return nil, fmt.Errorf("reading migration file %s: %w", filename, err)
		}
	}

	// Just to be sure
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Timestamp < changes[j].Timestamp
	})

	return changes, nil
}

func parseMigration(migrationStr []byte, filename string) (*Migration, error) {
	m := &Migration{
		Filename: filename,
	}

	var err error

	timestampString := timestampReg.FindString(filepath.Base(filename))
	m.Timestamp, err = strconv.Atoi(timestampString)
	if err != nil {
		return nil, err
	}

	m.Up, m.Down = parseChange(migrationStr)

	return m, nil
}

type section int

const (
	sectionNone section = iota
	sectionUp
	sectionDown
)

func parseChange(data []byte) (string, string) {
	var up, down strings.Builder

	var section section
	for line := range bytes.SplitSeq(data, []byte{'\n'}) {
		switch string(line) {
		case "-- migrate up":
			section = sectionUp
		case "-- migrate down":
			section = sectionDown
		default:
			switch section {
			case sectionUp:
				up.WriteByte('\n')
				up.Write(line)
			case sectionDown:
				down.WriteByte('\n')
				down.Write(line)
			case sectionNone:
			}
		}
	}

	return up.String(), down.String()
}
