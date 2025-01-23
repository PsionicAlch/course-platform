package vanillatext

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
)

type TemplateCache map[string]*template.Template

type Templates struct {
	Name  string
	Cache TemplateCache
}

// CreateTemplates constructs a new Templates instance based off the provided parameters.
func CreateTemplates(cdnURL string, embeddedFS embed.FS, directory string, otherDirectories ...string) (*Templates, error) {
	tmpls, err := embeddedFS.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	cache := make(TemplateCache, len(tmpls))

	for _, tmpl := range tmpls {
		if tmpl.IsDir() {
			continue
		}

		name := tmpl.Name()

		patterns := append(otherDirectories, fmt.Sprintf("%s/%s", directory, name))

		t, err := template.New(name).Funcs(CreateFuncMap(cdnURL)).ParseFS(embeddedFS, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = t
	}

	templates := &Templates{
		Name:  directory,
		Cache: cache,
	}

	return templates, nil
}

// Compile renders a given template to a buffer.
func (templates *Templates) Compile(tmpl string, data any) (*bytes.Buffer, error) {
	t, ok := templates.Cache[tmpl]
	if !ok {
		return nil, render.CreateFailedToFindTemplateInCache(tmpl, templates.Name)
	}

	templateBuffer := new(bytes.Buffer)

	err := t.Execute(templateBuffer, data)
	if err != nil {
		return nil, render.CreateFailedToCompileTemplate(tmpl, err)
	}

	return templateBuffer, nil
}
