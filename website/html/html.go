package html

import "embed"

//go:embed components/*.component.tmpl emails/*.email.tmpl htmx/*.htmx.tmpl layouts/*.layout.tmpl pages/*.page.tmpl rss/*.rss.tmpl
var HTMLFiles embed.FS
