package content

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database"
	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

type Content struct {
	utils.Loggers
}

func RegisterContent(db database.Database) error {
	loggers := utils.CreateLoggers("CONTENT")
	content := &Content{
		Loggers: loggers,
	}

	content.RegisterTutorialsContent(db)
	content.RegisterCourseContent(db)

	return nil
}
