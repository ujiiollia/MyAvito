package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// Banner представляет структуру баннера
type Banner struct {
	ID              int
	FeatureID       int
	Tags            []int
	JSONData        string
	Active          int
	LastUpdatedTime string
}

func New(storagePath string) (*Storage, error) {
	const el = "sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS banners (
		id INTEGER PRIMARY KEY,
		feature_id INTEGER,
		tag_id INTEGER,
		json_data TEXT,
		active INTEGER,
		last_updated_time DATETIME
	);
	`) //todo: CREATE INDEX IF NOT EXISTS

	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}

	return &Storage{db: db}, nil
}
