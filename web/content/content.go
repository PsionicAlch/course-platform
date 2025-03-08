package content

import (
	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/utils"
)

type Content struct {
	utils.Loggers
}

func RegisterContent(db database.Database) {
	loggers := utils.CreateLoggers("CONTENT")
	content := &Content{
		Loggers: loggers,
	}

	if err := db.DeleteAllKeywords(); err != nil {
		loggers.ErrorLog.Fatalf("Failed to delete all keywords: %s\n", err)
	}

	content.RegisterTutorialsContent(db)
	content.RegisterCourseContent(db)
}
