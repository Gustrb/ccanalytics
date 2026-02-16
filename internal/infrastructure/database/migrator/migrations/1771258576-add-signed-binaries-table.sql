-- migrate up
CREATE TABLE signed_binaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hash TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE UNIQUE INDEX idx_signed_binaries_hash ON signed_binaries (hash);

-- migrate down
DROP TABLE signed_binaries;