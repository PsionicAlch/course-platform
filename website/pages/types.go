package pages

import "github.com/PsionicAlch/psionicalch-home/internal/render"

type Renderers struct {
	Page render.Renderer
	Htmx render.Renderer
	RSS  render.Renderer
}

func CreateRenderers(pageRenderer render.Renderer, htmxRenderer render.Renderer, rssRenderer render.Renderer) *Renderers {
	return &Renderers{
		Page: pageRenderer,
		Htmx: htmxRenderer,
		RSS:  rssRenderer,
	}
}
