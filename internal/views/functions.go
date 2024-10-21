package views

import (
	"fmt"
	"html/template"
)

func CreateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"props": Props,
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
