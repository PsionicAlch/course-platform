package render

import (
	"errors"
	"fmt"
)

var (
	// ErrCannotRenderHTML represents the error status which states that the function called in unable to render any
	// HTML.
	ErrCannotRenderHTML = errors.New("renderer does not support rendering html")

	// ErrCannotRenderXML represents the error status which states that the function called in unable to render any
	// XML.
	ErrCannotRenderXML = errors.New("renderer does not support rendering xml")
)

// ErrFailedToFindTemplateInCache represents the error state in which the renderer could not find the requested
// template.
type ErrFailedToFindTemplateInCache error

// CreateFailedToFindTemplateInCache creates a new instance of ErrFailedToFindTemplateInCache based off the
// template and the name of template collection.
func CreateFailedToFindTemplateInCache(tmpl, templateName string) ErrFailedToFindTemplateInCache {
	return fmt.Errorf("failed to find \"%s\" in \"%s\"", tmpl, templateName)
}

// ErrFailedToCompileTemplate represents the error state in which the template failed to compile.
type ErrFailedToCompileTemplate error

// CreateFailedToCompileTemplate creates a new instance of ErrFailedToCompileTemplate based off the template and
// an error as to why the template could not be compiled.
func CreateFailedToCompileTemplate(tmpl string, err error) ErrFailedToCompileTemplate {
	return fmt.Errorf("failed to compile %s: %w", tmpl, err)
}
