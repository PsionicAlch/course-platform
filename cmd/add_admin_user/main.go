package main

import (
	"flag"

	"github.com/PsionicAlch/psionicalch-home/internal/authentication"
	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

func main() {
	loggers := utils.CreateLoggers("ADMIN USER COMMAND")

	db, err := sqlite_database.CreateSQLiteDatabase("/db/db.sqlite", "/db/migrations")
	if err != nil {
		loggers.ErrorLog.Fatalf("Failed to open database connection: %s\n", err)
	}
	defer db.Close()

	auth, err := authentication.SetupAuthentication(db, nil, 0, 0, "", "", "", "")
	if err != nil {
		loggers.ErrorLog.Fatalf("Failed to set up authentication: %s\n", err)
	}

	name := flag.String("name", "", "First name of the admin user.")
	surname := flag.String("surname", "", "Last name of the admin user.")
	email := flag.String("email", "", "Email address of the admin user.")
	password := flag.String("password", "", "Password of the admin user.")

	flag.Parse()

	if *name == "" || *surname == "" || *email == "" || *password == "" {
		loggers.ErrorLog.Println("To add a new admin user you need to specifically call:\n\tmake new-admin name=\"FIRST_NAME\" surname=\"LAST_NAME\" email=\"EMAIL_ADDRESS\" password=\"PASSWORD\"")
	}

	if err := auth.NewAdminUser(*name, *surname, *email, *password); err != nil {
		loggers.ErrorLog.Fatalf("Failed to add new admin user: %s\n", err)
	} else {
		loggers.InfoLog.Printf("%s %s has been added as an admin user!", *name, *surname)
	}
}
