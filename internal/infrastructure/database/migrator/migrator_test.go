package migrator

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	//go:embed migrations/1771180375-add-migrations-table.sql
	addMigrationsMigration string
)

func TestShouldBeAbleToParseMigrateUpFile(t *testing.T) {
	migration, err := parseMigration([]byte(addMigrationsMigration), "migrations/1771180375-add-migrations-table.sql")
	require.NoError(t, err)

	migrateUp := `CREATE TABLE migrations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	"timestamp" INTEGER NOT NULL,
	"filename" TEXT NOT NULL,
    "hash" TEXT NOT NULL,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);`
	require.Equal(t, simplifyString(migrateUp), simplifyString(migration.Up))
}

func simplifyString(s string) string {
	// remove \n and \t
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}
