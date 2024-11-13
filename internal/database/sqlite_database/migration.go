package sqlite_database

import (
	"github.com/golang-migrate/migrate/v4"
)

func (db *SQLiteDatabase) MigrateUp() error {
	m, err := db.SetupMigrations()
	if err != nil {
		return err
	}

	// Apply all up migrations.
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		db.ErrorLog.Printf("Failed to run up migrations: %s", err.Error())
		return err
	}

	return nil
}

func (db *SQLiteDatabase) MigrateDown() error {
	m, err := db.SetupMigrations()
	if err != nil {
		return err
	}

	// Apply all down migrations.
	if err = m.Down(); err != nil && err != migrate.ErrNoChange {
		db.ErrorLog.Printf("Failed to run down migrations: %s", err.Error())
		return err
	}

	return nil
}

func (db *SQLiteDatabase) Rollback(steps int) error {
	m, err := db.SetupMigrations()
	if err != nil {
		return err
	}

	if err = m.Steps(-1 * steps); err != nil && err != migrate.ErrNoChange {
		db.ErrorLog.Printf("Failed to rollback migrations: %s", err.Error())
		return err
	}

	return nil
}
