package vanillatext

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

func CreateFuncMap(cdnURL string) template.FuncMap {
	funcMap := template.FuncMap{
		"current_time":           CurrentTime,
		"format_time_to_rfc_822": FormatTimeToRFC822,
		"assets":                 Assets(cdnURL),
	}

	return funcMap
}

func CurrentTime() time.Time {
	return time.Now()
}

func FormatTimeToRFC822(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

func Assets(cdnURL string) func(path string) string {
	return func(path string) string {
		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}

		return cdnURL + path
	}
}
