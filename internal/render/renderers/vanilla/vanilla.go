package vanilla

import (
	"embed"
	"io"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
)

type VanillaRenderer struct {
	templates *Templates
}

func SetupVanillaRenderer(embeddedFS embed.FS, directory string, otherDirectories ...string) (*VanillaRenderer, error) {
	templates, err := CreateTemplates(embeddedFS, directory, otherDirectories...)
	if err != nil {
		return nil, err
	}

	vanillaRenderer := &VanillaRenderer{
		templates: templates,
	}

	return vanillaRenderer, nil
}

func (renderer *VanillaRenderer) Render(w io.Writer, tmpl string, data any) error {
	tmplBuffer, err := renderer.templates.Compile(tmpl, data)
	if err != nil {
		return err
	}

	_, err = tmplBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func (renderer *VanillaRenderer) RenderHTML(w http.ResponseWriter, tmpl string, data any, status ...int) error {
	tmplBuffer, err := renderer.templates.Compile(tmpl, data)
	if err != nil {
		return err
	}

	statusCode := render.GetStatusCode(status...)
	w.WriteHeader(statusCode)

	_, err = tmplBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
