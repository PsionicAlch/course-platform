package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

type Loggers struct {
	InfoLog    *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
}

// CreateLoggers creates an info, warning and error logger based on the provided prefix.
func CreateLoggers(prefix string) Loggers {
	infoPrefix := color.CyanString(fmt.Sprintf("%s INFO --->\t", prefix))
	warningPrefix := color.YellowString(fmt.Sprintf("%s WARNING --->\t", prefix))
	errorPrefix := color.RedString(fmt.Sprintf("%s ERROR --->\t", prefix))

	return Loggers{
		InfoLog:    log.New(os.Stdout, infoPrefix, log.Ltime|log.Ldate),
		WarningLog: log.New(os.Stdout, warningPrefix, log.Ltime|log.Ldate|log.Lshortfile),
		ErrorLog:   log.New(os.Stderr, errorPrefix, log.Ltime|log.Ldate|log.Lshortfile),
	}
}
