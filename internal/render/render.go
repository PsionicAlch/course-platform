package render

import (
	"context"
	"io"
	"net/http"
)

// Renderer represents all functions available from the render package.
type Renderer interface {
	Render(w io.Writer, ctx context.Context, file string, data any) error
	RenderHTML(w http.ResponseWriter, ctx context.Context, file string, data any, status ...int) error
	RenderXML(w http.ResponseWriter, file string, data any, status ...int) error
}
