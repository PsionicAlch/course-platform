package vanilla

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
)

func CreateFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"props":    Props,
		"keywords": Keywords,
		"html":     HTML,
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

func Keywords(tutorial *models.TutorialModel) string {
	var keywords []string
	for _, keywordModel := range tutorial.Keywords {
		keywords = append(keywords, keywordModel.Keyword)
	}

	return strings.Join(keywords, ", ")
}

func HTML(s string) template.HTML {
	return template.HTML(s)
}
