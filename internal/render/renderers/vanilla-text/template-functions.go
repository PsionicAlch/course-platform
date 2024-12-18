package vanillatext

import (
	"text/template"
	"time"
)

func CreateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"current_time":           CurrentTime,
		"format_time_to_rfc_822": FormatTimeToRFC822,
	}

	return funcMap
}

func CurrentTime() time.Time {
	return time.Now()
}

func FormatTimeToRFC822(t time.Time) string {
	return t.Format(time.RFC1123Z)
}
