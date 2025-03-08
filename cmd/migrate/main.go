package main

import (
	"os"
	"strconv"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	"github.com/PsionicAlch/course-platform/internal/database/sqlite_database"
	"github.com/PsionicAlch/course-platform/internal/utils"
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

	db, err := sqlite_database.CreateSQLiteDatabase("/db/db.sqlite", "/db/migrations")
	if err != nil {
		loggers.ErrorLog.Fatalln(err)
	}

	// Run command.
	switch subcommand {
	case "up":
		loggers.InfoLog.Println("Running migrations up.")

		err = db.MigrateUp()
		if err != nil {
			loggers.ErrorLog.Fatalln(err)
		}

		loggers.InfoLog.Println("Migrations applied successfully!")

	case "down":
		loggers.InfoLog.Println("Running migrations down.")

		err = db.MigrateDown()
		if err != nil {
			loggers.ErrorLog.Fatalln(err)
		}

		loggers.InfoLog.Println("Migrations removed successfully!")

	case "rollback":
		// Parse the optional `-step` flag
		steps := 1

		if len(os.Args) > 2 && os.Args[2] == "-step" && len(os.Args) > 3 {
			steps, _ = strconv.Atoi(os.Args[3]) // Convert string to int
		}

		loggers.InfoLog.Printf("Rolling back %d step(s)\n", steps)

		err = db.Rollback(steps)
		if err != nil {
			loggers.ErrorLog.Fatalln(err)
		}

		loggers.InfoLog.Printf("Successfully rolled back %d migration(s)\n", steps)

	default:
		loggers.ErrorLog.Fatalf("unknown subcommand: %s", subcommand)
	}
}
