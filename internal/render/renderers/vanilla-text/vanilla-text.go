package vanillatext

import (
	"context"
	"embed"
	"io"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
)

type VanillaTextRenderer struct {
	templates     *Templates
	fileExtension string
}

func SetupVanillaTextRenderer(embeddedFS embed.FS, fileExtension, directory string, otherDirectories ...string) (*VanillaTextRenderer, error) {
	templates, err := CreateTemplates(embeddedFS, directory, otherDirectories...)
	if err != nil {
		return nil, err
	}

	vanillaTextRenderer := &VanillaTextRenderer{
		templates:     templates,
		fileExtension: fileExtension,
	}

	return vanillaTextRenderer, nil
}

func (renderer *VanillaTextRenderer) Render(w io.Writer, ctx context.Context, file string, data any) error {
	tmplBuffer, err := renderer.templates.Compile(file+renderer.fileExtension, data)
	if err != nil {
		return err
	}

	_, err = tmplBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func (renderer *VanillaTextRenderer) RenderHTML(w http.ResponseWriter, ctx context.Context, file string, data any, status ...int) error {
	return render.ErrCannotRenderHTML
}

func (renderer *VanillaTextRenderer) RenderXML(w http.ResponseWriter, file string, data any, status ...int) error {
	tmplBuffer, err := renderer.templates.Compile(file+renderer.fileExtension, data)
	if err != nil {
		return err
	}

	statusCode := render.GetStatusCode(status...)
	w.WriteHeader(statusCode)

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")

	_, err = tmplBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
