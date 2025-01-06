package vanillahtml

import (
	"fmt"
	"html/template"
	"net/url"
	"time"
)

func CreateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"props":                   Props,
		"pretty_date":             PrettyDate,
		"format_time_to_iso_8601": FormatTimeToISO8601,
		"html":                    HTML,
		"add_queries":             AddQueries,
		"current_time":            CurrentTime,
		"url_escape":              URLEscape,
	}

	return funcMap
}

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

func PrettyDate(t time.Time) string {
	return t.Format("Monday, January 2, 2006 at 3:04 PM")
}

func FormatTimeToISO8601(t time.Time) string {
	return t.Format("2006-01-02")
}

func HTML(s string) template.HTML {
	return template.HTML(s)
}

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

func CurrentTime() time.Time {
	return time.Now()
}

func URLEscape(s string) string {
	return url.QueryEscape(s)
}
