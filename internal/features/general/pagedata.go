package general

import "github.com/PsionicAlch/psionicalch-home/internal/views"

type HomePageData struct {
	HeaderComponentData *views.HeaderComponentData
}

func CreateHomePageData() *HomePageData {
	headerComponentData := views.CreateHeaderComponent(
		"Build Real Projects, Launch Real Products",
		"Tired of courses that teach theory but don't get real-world results? Skip the fluff and dive into practical Golang courses designed to help you build projects you can use for your portfolio or business right now.",
		"/courses",
		"Start Building Today",
	)

	return &HomePageData{
		HeaderComponentData: headerComponentData,
	}
}
