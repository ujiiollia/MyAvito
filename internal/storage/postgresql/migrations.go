package postge

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrater interface {
	Up() error
	Close() (sourceErr, databaseErr error)
}

func ApplyMigrations(m Migrater) error {
	const el = "postgresql.migrations.ApplyMigrations"

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s: %w", el, err)
	}

	sourceErr, databaseErr := m.Close()

	if sourceErr != nil {
		return fmt.Errorf("%s: %w", el, sourceErr)
	}

	if databaseErr != nil {
		return fmt.Errorf("%s: %w", el, databaseErr)
	}

	return nil
}
