package html

import "embed"

//go:embed components/*.component.tmpl emails/*.email.tmpl htmx/*.htmx.tmpl layouts/*.layout.tmpl pages/*.page.tmpl
var HTMLFiles embed.FS

//go:embed rss/*.rss.tmpl
var TextFiles embed.FS
