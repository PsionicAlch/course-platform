package main

import (
	"time"

	"github.com/PsionicAlch/course-platform/internal/database/sqlite_database"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/content"
)

func main() {
	startTimer := time.Now()

	loggers := utils.CreateLoggers("CONTENT LOADER")

	loggers.InfoLog.Println("Creating database connection!")

	db, err := sqlite_database.CreateSQLiteDatabase("/db/db.sqlite", "/db/migrations")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to open database connection: ", err)
	}
	defer db.Close()

	loggers.InfoLog.Println("Registering content!")

	content.RegisterContent(db)

	endTimer := time.Since(startTimer)

	loggers.InfoLog.Printf("Finished loading content in %s!", endTimer)
}
