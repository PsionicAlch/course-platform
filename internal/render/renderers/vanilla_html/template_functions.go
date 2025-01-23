package vanillahtml

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"
)

// CreateFuncMap constructs a function map for all HTML based functions.
func CreateFuncMap(cdnURL string) template.FuncMap {
	funcMap := template.FuncMap{
		"props":                   Props,
		"pretty_date":             PrettyDate,
		"format_time_to_iso_8601": FormatTimeToISO8601,
		"html":                    HTML,
		"add_queries":             AddQueries,
		"current_time":            CurrentTime,
		"url_escape":              URLEscape,
		"assets":                  Assets(cdnURL),
	}

	return funcMap
}

// Props constructs a map that can be passed to HTML components as properties.
func Props(values ...any) (map[string]any, error) {
	valuesLen := len(values)

	if valuesLen%2 != 0 {
		return nil, fmt.Errorf("dict requires an even number of arguments")
	}

	dict := make(map[string]interface{}, valuesLen/2)
	for i := 0; i < valuesLen; i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}

	return dict, nil
}

// PrettyDate formats a date into something more human readable.
func PrettyDate(t time.Time) string {
	return t.Format("Monday, January 2, 2006 at 3:04 PM")
}

// FormatTimeToISO8601 formats a date according to the ISO 8601 standard.
func FormatTimeToISO8601(t time.Time) string {
	return t.Format("2006-01-02")
}

// HTML converts a given string to a template.HTML type so that HTML can be rendered without being
// escaped.
func HTML(s string) template.HTML {
	return template.HTML(s)
}

// AddQueries adds values to a URL query.
func AddQueries(queries string, values ...any) (string, error) {
	valuesLen := len(values)

	if valuesLen%2 != 0 {
		return "", fmt.Errorf("url queries need to be an even number (Query key/value pair)")
	}

	urlQuery, err := url.ParseQuery(queries)
	if err != nil {
		return "", fmt.Errorf("failed to parse url queries: %s", err)
	}

	for i := 0; i < valuesLen; i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return "", fmt.Errorf("url query keys must be strings")
		}

		value := fmt.Sprintf("%v", values[i+1])

		if urlQuery.Has(key) {
			urlQuery.Set(key, value)
		} else {
			urlQuery.Add(key, value)
		}
	}

	return urlQuery.Encode(), nil
}

// CurrentTime returns the current time.
func CurrentTime() time.Time {
	return time.Now()
}

// URLEscape ensures that a string is URL safe.
func URLEscape(s string) string {
	return url.QueryEscape(s)
}

// Assets appends a path to the provided CDN link.
func Assets(cdnURL string) func(path string) string {
	return func(path string) string {
		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}

		return cdnURL + path
	}
}
