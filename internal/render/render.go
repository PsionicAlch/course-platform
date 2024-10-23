package render

import (
	"io"
	"net/http"
)

type Renderer interface {
	Render(w io.Writer, tmpl string, data any) error
	RenderHTML(w http.ResponseWriter, tmpl string, data any, status ...int) error
}
