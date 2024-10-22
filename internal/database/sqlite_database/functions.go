package sqlite_database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func (db *SQLiteDatabase) SetupMigrations() (*migrate.Migrate, error) {
	// Configure database driver for migrate.
	driver, err := sqlite.WithInstance(db.connection, &sqlite.Config{})
	if err != nil {
		return nil, err
	}

	// Specify migrations directory.
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file:/%s", db.migrationsDir), "sqlite", driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}
