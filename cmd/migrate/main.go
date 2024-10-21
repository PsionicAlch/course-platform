package main

import (
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

func main() {
	// Construct loggers.
	loggers := utils.CreateLoggers("MIGRATIONS")

	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		loggers.ErrorLog.Fatal("expected 'up', 'down' or 'rollback' subcommand")
	}

	// Get the subcommand (e.g., "up", "down", "rollback")
	subcommand := os.Args[1]

	// Open database.
	db := database.CreateDatabase()
	defer db.Close()

	// Configure database driver for migrate.
	driver, err := sqlite.WithInstance(db.GetConnection(), &sqlite.Config{})
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to configure database driver for migrate: ", err)
	}

	// Specify migrations directory.
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "sqlite", driver)
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to construct new instance of migrate: ", err)
	}

	// Run command.
	switch subcommand {
	case "up":
		loggers.InfoLog.Println("Running migrations up.")

		// Apply all up migrations.
		if err = m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				loggers.InfoLog.Println("Migrations are up to date.")
			} else {
				loggers.ErrorLog.Fatalln("Failed to run migrations: ", err)
			}
		}

		loggers.InfoLog.Println("Migrations applied successfully!")

	case "down":
		loggers.InfoLog.Println("Running migrations down.")

		// Apply all down migrations.
		if err = m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				loggers.InfoLog.Println("Migrations are up to date.")
			} else {
				loggers.ErrorLog.Fatalln("Failed to run migrations: ", err)
			}
		}

		loggers.InfoLog.Println("Migrations removed successfully!")

	case "rollback":
		// Parse the optional `-step` flag
		steps := 1

		if len(os.Args) > 2 && os.Args[2] == "-step" && len(os.Args) > 3 {
			steps, _ = strconv.Atoi(os.Args[3]) // Convert string to int
		}

		loggers.InfoLog.Printf("Rolling back %d step(s)\n", steps)

		if err = m.Steps(-1 * steps); err != nil {
			if err == migrate.ErrNoChange {
				loggers.InfoLog.Println("Migrations are up to date.")
			} else {
				loggers.ErrorLog.Fatalln("Failed to run migrations: ", err)
			}
		}

		loggers.InfoLog.Printf("Successfully rolled back %d migration(s)\n", steps)

	default:
		loggers.ErrorLog.Fatalf("unknown subcommand: %s", subcommand)
	}
}
