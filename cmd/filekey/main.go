package main

import (
	"log"

	"github.com/PsionicAlch/psionicalch-home/internal/database/sqlite_database"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/PsionicAlch/psionicalch-home/web/content"
)

func main() {
	loggers := utils.CreateLoggers("FILE KEY GENERATOR")

	db, err := sqlite_database.CreateSQLiteDatabase("/db/db.sqlite", "/db/migrations")
	if err != nil {
		loggers.ErrorLog.Fatalln("Failed to open database connection: ", err)
	}
	defer db.Close()

	for {
		key, err := content.GenerateFileKey()
		if err != nil {
			log.Fatalln(err)
		}

		tutorial, err := db.GetTutorialByFileKey(key)
		if err != nil {
			loggers.ErrorLog.Fatalf("Failed to try and get tutorial by file key: %s\n", err)
		}

		course, err := db.GetCourseByFileKey(key)
		if err != nil {
			loggers.ErrorLog.Fatalf("Failed to try and get course by file key: %s\n", err)
		}

		chapter, err := db.GetChapterByFileKey(key)
		if err != nil {
			loggers.ErrorLog.Fatalf("Failed to try and get chapter by file key: %s\n", err)
		}

		if tutorial == nil && course == nil && chapter == nil {
			loggers.InfoLog.Printf("File Key: %s\n", key)
			break
		}
	}
}
