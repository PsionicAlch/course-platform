package views

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
)

//go:embed components htmx layouts pages
var viewsFS embed.FS

type TemplateCache map[string]*template.Template

type Views struct {
	utils.Loggers
	pagesCache TemplateCache
	htmxCache  TemplateCache
}

func SetupRenderer() *Views {
	renderLoggers := utils.CreateLoggers("VIEWS")

	pagesTemplateCache, err := CreatePagesTemplateCache()
	if err != nil {
		renderLoggers.ErrorLog.Fatal("Failed to construct pages template cache: ", err)
	}

	htmxTemplateCache, err := CreateHTMXTemplateCache()
	if err != nil {
		renderLoggers.ErrorLog.Fatal("Failed to construct htmx template cache: ", err)
	}

	return &Views{
		Loggers:    renderLoggers,
		pagesCache: pagesTemplateCache,
		htmxCache:  htmxTemplateCache,
	}
}

func CreatePagesTemplateCache() (TemplateCache, error) {
	pages, err := viewsFS.ReadDir("pages")
	if err != nil {
		return nil, err
	}

	templateCache := make(map[string]*template.Template, len(pages))

	for _, page := range pages {
		if page.IsDir() {
			continue
		}

		name := page.Name()

		t, err := template.New(name).Funcs(CreateFuncMap()).ParseFS(viewsFS, "components/*.component.tmpl", "layouts/*.layout.tmpl", "pages/"+name)
		if err != nil {
			return nil, err
		}

		templateCache[name] = t
	}

	return templateCache, nil
}

func CreateHTMXTemplateCache() (TemplateCache, error) {
	components, err := viewsFS.ReadDir("htmx")
	if err != nil {
		return nil, err
	}

	templateCache := make(map[string]*template.Template, len(components))

	for _, component := range components {
		if component.IsDir() {
			continue
		}

		name := component.Name()

		t, err := template.New(name).Funcs(CreateFuncMap()).ParseFS(viewsFS, "components/*.component.tmpl", "htmx/"+name)
		if err != nil {
			return nil, err
		}

		templateCache[name] = t
	}

	return templateCache, nil
}

func (v *Views) RenderHTML(w http.ResponseWriter, tmpl string, data any, status ...int) {
	var statusCode int

	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = http.StatusOK
	}

	t, ok := v.pagesCache[tmpl]
	if !ok {
		v.ErrorLog.Printf("Failed to find %s in pages template cache!\n", tmpl)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	err := t.Execute(w, data)
	if err != nil {
		v.ErrorLog.Println(err)
		return
	}
}

func (v *Views) RenderHTMX(w http.ResponseWriter, tmpl string, data any, status ...int) {
	var statusCode int

	if len(status) > 0 {
		statusCode = status[0]
	} else {
		statusCode = http.StatusOK
	}

	t, ok := v.htmxCache[tmpl]
	if !ok {
		v.ErrorLog.Printf("Failed to find %s in htmx template cache!\n", tmpl)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	err := t.Execute(w, data)
	if err != nil {
		v.ErrorLog.Println(err)
		return
	}
}

func (v *Views) RenderNotFound(w http.ResponseWriter) {
	v.RenderHTML(w, "404.page.tmpl", nil, http.StatusNotFound)
}

func (v *Views) RenderInternalServerError(w http.ResponseWriter) {
	v.RenderHTML(w, "500.page.tmpl", nil, http.StatusInternalServerError)
}
