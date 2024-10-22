package views

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/PsionicAlch/psionicalch-home/internal/views/errors"
)

//go:embed components htmx layouts pages
var viewsFS embed.FS

// type TemplateCache map[string]*template.Template
type TemplateCache struct {
	Name  string
	Cache map[string]*template.Template
}

type Views struct {
	pagesCache *TemplateCache
	htmxCache  *TemplateCache
}

func SetupRenderer() (*Views, error) {
	pagesTemplateCache, err := CreateTemplateCache("pages", "components/*.component.tmpl", "layouts/*.layout.tmpl")
	if err != nil {
		return nil, errors.CreateFailedToConstructTemplateCacheError(fmt.Sprintf("failed to construct pages template cache: %s", err.Error()))
	}

	htmxTemplateCache, err := CreateTemplateCache("htmx", "components/*.component.tmpl")
	if err != nil {
		return nil, errors.CreateFailedToConstructTemplateCacheError(fmt.Sprintf("failed to construct htmx template cache: %s", err.Error()))
	}

	v := &Views{
		pagesCache: pagesTemplateCache,
		htmxCache:  htmxTemplateCache,
	}

	return v, nil
}

func CreateTemplateCache(directory string, otherDirs ...string) (*TemplateCache, error) {
	tmpls, err := viewsFS.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	cache := make(map[string]*template.Template, len(tmpls))

	for _, tmpl := range tmpls {
		if tmpl.IsDir() {
			continue
		}

		name := tmpl.Name()

		patterns := append(otherDirs, fmt.Sprintf("%s/%s", directory, name))
		t, err := template.New(name).Funcs(CreateFuncMap()).ParseFS(viewsFS, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = t
	}

	templateCache := &TemplateCache{
		Name:  directory,
		Cache: cache,
	}

	return templateCache, nil
}

// func CreatePagesTemplateCache() (TemplateCache, error) {
// 	pages, err := viewsFS.ReadDir("pages")
// 	if err != nil {
// 		return nil, err
// 	}

// 	templateCache := make(map[string]*template.Template, len(pages))

// 	for _, page := range pages {
// 		if page.IsDir() {
// 			continue
// 		}

// 		name := page.Name()

// 		t, err := template.New(name).Funcs(CreateFuncMap()).ParseFS(viewsFS, "components/*.component.tmpl", "layouts/*.layout.tmpl", "pages/"+name)
// 		if err != nil {
// 			return nil, err
// 		}

// 		templateCache[name] = t
// 	}

// 	return templateCache, nil
// }

// func CreateHTMXTemplateCache() (TemplateCache, error) {
// 	components, err := viewsFS.ReadDir("htmx")
// 	if err != nil {
// 		return nil, err
// 	}

// 	templateCache := make(map[string]*template.Template, len(components))

// 	for _, component := range components {
// 		if component.IsDir() {
// 			continue
// 		}

// 		name := component.Name()

// 		t, err := template.New(name).Funcs(CreateFuncMap()).ParseFS(viewsFS, "components/*.component.tmpl", "htmx/"+name)
// 		if err != nil {
// 			return nil, err
// 		}

// 		templateCache[name] = t
// 	}

// 	return templateCache, nil
// }

func (v *Views) RenderHTML(tmpl string, data any) (*bytes.Buffer, error) {
	return Render(v.pagesCache, tmpl, data)
}

func (v *Views) RenderHTMX(tmpl string, data any) (*bytes.Buffer, error) {
	return Render(v.htmxCache, tmpl, data)
}

// func (v *Views) RenderHTML(w http.ResponseWriter, tmpl string, data any, status ...int) error {
// 	var statusCode int

// 	if len(status) > 0 {
// 		statusCode = status[0]
// 	} else {
// 		statusCode = http.StatusOK
// 	}

// 	t, ok := v.pagesCache[tmpl]
// 	if !ok {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return errors.CreateFailedToFindTemplateInCache(tmpl, "pages")
// 	}

// 	w.WriteHeader(statusCode)

// 	err := t.Execute(w, data)
// 	if err != nil {
// 		return errors.CreateFailedToCompileTemplate(tmpl, err.Error())
// 	}

// 	return nil
// }

// func (v *Views) RenderHTMX(w http.ResponseWriter, tmpl string, data any, status ...int) {
// 	var statusCode int

// 	if len(status) > 0 {
// 		statusCode = status[0]
// 	} else {
// 		statusCode = http.StatusOK
// 	}

// 	t, ok := v.htmxCache[tmpl]
// 	if !ok {
// 		v.ErrorLog.Printf("Failed to find %s in htmx template cache!\n", tmpl)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(statusCode)

// 	err := t.Execute(w, data)
// 	if err != nil {
// 		v.ErrorLog.Println(err)
// 		return
// 	}
// }

// func (v *Views) RenderNotFound(w http.ResponseWriter) {
// 	v.RenderHTML(w, "404.page.tmpl", nil, http.StatusNotFound)
// }

// func (v *Views) RenderInternalServerError(w http.ResponseWriter) {
// 	v.RenderHTML(w, "500.page.tmpl", nil, http.StatusInternalServerError)
// }
