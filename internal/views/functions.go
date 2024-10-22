package views

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/PsionicAlch/psionicalch-home/internal/views/errors"
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

func Render(templateCache *TemplateCache, tmpl string, data any) (*bytes.Buffer, error) {
	t, ok := templateCache.Cache[tmpl]
	if !ok {
		return nil, errors.CreateFailedToFindTemplateInCache(tmpl, templateCache.Name)
	}

	templateBuffer := new(bytes.Buffer)

	err := t.Execute(templateBuffer, data)
	if err != nil {
		return nil, errors.CreateFailedToCompileTemplate(tmpl, err.Error())
	}

	return templateBuffer, nil
}
