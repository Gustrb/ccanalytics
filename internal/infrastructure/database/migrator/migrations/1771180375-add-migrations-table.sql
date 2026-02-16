-- migrate up
CREATE TABLE migrations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	"timestamp" INTEGER NOT NULL,
	"filename" TEXT NOT NULL,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL
);

-- migrate down
DROP TABLE migrations;