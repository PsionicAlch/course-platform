package content

import (
	"sync"

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

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go content.RegisterTutorialsContent(&waitGroup, db)

	waitGroup.Wait()

	return nil
}
