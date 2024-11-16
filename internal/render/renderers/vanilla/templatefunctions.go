package vanilla

import (
	"fmt"
	"html/template"
	"time"
)

func CreateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"props":       Props,
		"pretty_date": PrettyDate,
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
