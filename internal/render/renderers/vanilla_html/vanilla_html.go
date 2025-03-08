package vanillahtml

import (
	"context"
	"embed"
	"io"
	"net/http"

	"github.com/PsionicAlch/course-platform/internal/render"
	"github.com/PsionicAlch/course-platform/internal/session"
)

type VanillaHTMLRenderer struct {
	templates     *Templates
	fileExtension string
	sessions      *session.Session
}

// SetupVanillaHTMLRenderer creates a new instance of the VanillaHTMLRenderer based on the provided properties.
func SetupVanillaHTMLRenderer(cdnURL string, sessions *session.Session, embeddedFS embed.FS, fileExtension, directory string, otherDirectories ...string) (*VanillaHTMLRenderer, error) {
	templates, err := CreateTemplates(cdnURL, embeddedFS, directory, otherDirectories...)
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

// Render writes the compiled template to the provided io.Writer.
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

// RenderHTML writes the compiled directly to the http.ResponseWriter as well as sets the necessary headers.
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

// RenderXML should not be called and will only ever return an error.
func (renderer *VanillaHTMLRenderer) RenderXML(w http.ResponseWriter, file string, data any, status ...int) error {
	return render.ErrCannotRenderXML
}
