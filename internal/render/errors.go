package render

import (
	"errors"
	"fmt"
)

var (
	ErrCannotRenderHTML = errors.New("renderer does not support rendering html")
	ErrCannotRenderXML  = errors.New("renderer does not support rendering xml")
)

type ErrFailedToFindTemplateInCache error

func CreateFailedToFindTemplateInCache(tmpl, templateName string) ErrFailedToFindTemplateInCache {
	return fmt.Errorf("failed to find \"%s\" in \"%s\"", tmpl, templateName)
}

type ErrFailedToCompileTemplate error

func CreateFailedToCompileTemplate(tmpl string, err error) ErrFailedToCompileTemplate {
	return fmt.Errorf("failed to compile %s: %w", tmpl, err)
}
