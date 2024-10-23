package html

import "embed"

//go:embed components/*.component.tmpl htmx/*.htmx.tmpl layouts/*.layout.tmpl pages/*.page.tmpl
var HTMLFiles embed.FS
