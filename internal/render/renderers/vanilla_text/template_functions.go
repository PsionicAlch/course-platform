package vanillatext

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

// CreateFuncMap constructs a map of functions that can be used in the text templates.
func CreateFuncMap(cdnURL string) template.FuncMap {
	funcMap := template.FuncMap{
		"current_time":           CurrentTime,
		"format_time_to_rfc_822": FormatTimeToRFC822,
		"assets":                 Assets(cdnURL),
	}

	return funcMap
}

// CurrentTime returns the current time.
func CurrentTime() time.Time {
	return time.Now()
}

// FormatTimeToRFC822 formats the time according the RFC 822 standards.
func FormatTimeToRFC822(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

// Assets constructs a URL based off the given path and the CDN url.
func Assets(cdnURL string) func(path string) string {
	return func(path string) string {
		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}

		return cdnURL + path
	}
}
