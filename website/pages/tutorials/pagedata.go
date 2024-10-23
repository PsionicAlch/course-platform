package tutorials

import "github.com/PsionicAlch/psionicalch-home/website/html"

type TutorialsPageData struct {
	HeaderComponentData *html.HeaderComponentData
	AllTutorials        []*TutorialData
}

func CreateTutorialsPageData(tutorials *Tutorials) *TutorialsPageData {
	headerComponentData := html.CreateHeaderComponent(
		"Quick Tutorials, Real-World Skills",
		"Our bite-sized tutorials give you the skills you need without the fluff. Each tutorial is a practical snippet from our in-depth courses, helping you build real-world projects one step at a time. Whether you're short on time or looking for focused learning, these tutorials get you up to speed fast.",
		"#tutorials",
		"Start Learning",
	)
	allTutorials := tutorials.GetAllTutorials()

	return &TutorialsPageData{
		HeaderComponentData: headerComponentData,
		AllTutorials:        allTutorials,
	}
}
