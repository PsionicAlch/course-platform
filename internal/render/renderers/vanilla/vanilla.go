package vanilla

import (
	"context"
	"embed"
	"io"
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/internal/render"
	"github.com/PsionicAlch/psionicalch-home/internal/session"
)

type VanillaRenderer struct {
	templates     *Templates
	fileExtension string
	sessions      *session.Session
}

func SetupVanillaRenderer(sessions *session.Session, embeddedFS embed.FS, fileExtension, directory string, otherDirectories ...string) (*VanillaRenderer, error) {
	templates, err := CreateTemplates(embeddedFS, directory, otherDirectories...)
	if err != nil {
		return nil, err
	}

	vanillaRenderer := &VanillaRenderer{
		templates:     templates,
		fileExtension: fileExtension,
		sessions:      sessions,
	}

	return vanillaRenderer, nil
}

func (renderer *VanillaRenderer) Render(w io.Writer, ctx context.Context, file string, data any) error {
	var infoMessages, warningMessages, errorMessages []string

	if renderer.sessions != nil {
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

func (renderer *VanillaRenderer) RenderHTML(w http.ResponseWriter, ctx context.Context, file string, data any, status ...int) error {
	var infoMessages, warningMessages, errorMessages []string

	if renderer.sessions != nil {
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

	_, err = tmplBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
