package vanillahtml

import (
	"context"
	"embed"
	"io"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
)

type VanillaHTMLRenderer struct {
	templates     *Templates
	fileExtension string
	sessions      *session.Session
}

func SetupVanillaHTMLRenderer(sessions *session.Session, embeddedFS embed.FS, fileExtension, directory string, otherDirectories ...string) (*VanillaHTMLRenderer, error) {
	templates, err := CreateTemplates(embeddedFS, directory, otherDirectories...)
	if err != nil {
		return nil, err
	}

	vanillaHTMLRenderer := &VanillaHTMLRenderer{
		templates:     templates,
		fileExtension: fileExtension,
		sessions:      sessions,
	}

	return vanillaHTMLRenderer, nil
}

func (renderer *VanillaHTMLRenderer) Render(w io.Writer, ctx context.Context, file string, data any) error {
	var infoMessages, warningMessages, errorMessages []string

	if renderer.sessions != nil && ctx != nil {
		infoMessages = renderer.sessions.GetInfoMessages(ctx)
		warningMessages = renderer.sessions.GetWarningMessages(ctx)
		errorMessages = renderer.sessions.GetErrorMessages(ctx)
	}

	dataMap := map[string]any{
		"UserData":        data,
		"InfoMessages":    infoMessages,
		"WarningMessages": warningMessages,
		"ErrorMessages":   errorMessages,
	}

	tmplBuffer, err := renderer.templates.Compile(file+renderer.fileExtension, dataMap)
	if err != nil {
		return err
	}

	_, err = tmplBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func (renderer *VanillaHTMLRenderer) RenderHTML(w http.ResponseWriter, ctx context.Context, file string, data any, status ...int) error {
	var infoMessages, warningMessages, errorMessages []string

	if renderer.sessions != nil && ctx != nil {
		infoMessages = renderer.sessions.GetInfoMessages(ctx)
		warningMessages = renderer.sessions.GetWarningMessages(ctx)
		errorMessages = renderer.sessions.GetErrorMessages(ctx)
	}

	dataMap := map[string]any{
		"UserData":        data,
		"InfoMessages":    infoMessages,
		"WarningMessages": warningMessages,
		"ErrorMessages":   errorMessages,
	}

	tmplBuffer, err := renderer.templates.Compile(file+renderer.fileExtension, dataMap)
	if err != nil {
		return err
	}

	statusCode := render.GetStatusCode(status...)
	w.WriteHeader(statusCode)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err = tmplBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func (renderer *VanillaHTMLRenderer) RenderXML(w http.ResponseWriter, file string, data any, status ...int) error {
	return render.ErrCannotRenderXML
}
