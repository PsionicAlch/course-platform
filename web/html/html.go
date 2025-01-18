package html

import "embed"

//go:embed components/*.component.tmpl emails/*.email.tmpl htmx/*.htmx.tmpl layouts/*.layout.tmpl pages/*.page.tmpl
var HTMLFiles embed.FS

//go:embed xml/*.xml.tmpl
var XMLFiles embed.FS

// TODO: Set up favicon to work with pages that don't use any layouts.
