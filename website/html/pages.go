package html

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/website/pages/tutorials"
)

type SignUpPageData struct {
	SignUpForm *SignupFormComponentData
}

type LoginPageData struct {
	LoginForm *LoginFormComponentData
}

type ProfilePageData struct {
	Email string
}

type HomePageData struct {
	HeaderComponentData *HeaderComponentData
}

func CreateHomePageData(user *models.UserModel) *HomePageData {
	headerComponentData := CreateHeaderComponent(
		"Build Real Projects, Launch Real Products",
		"Tired of courses that teach theory but don't get real-world results? Skip the fluff and dive into practical Golang courses designed to help you build projects you can use for your portfolio or business right now.",
		"/courses",
		"Start Building Today",
		user,
	)

	return &HomePageData{
		HeaderComponentData: headerComponentData,
	}
}

type TutorialsPageData struct {
	HeaderComponentData *HeaderComponentData
	AllTutorials        []*tutorials.TutorialData
}

func CreateTutorialsPageData(user *models.UserModel, tutorials *tutorials.Tutorials) *TutorialsPageData {
	headerComponentData := CreateHeaderComponent(
		"Quick Tutorials, Real-World Skills",
		"Our bite-sized tutorials give you the skills you need without the fluff. Each tutorial is a practical snippet from our in-depth courses, helping you build real-world projects one step at a time. Whether you're short on time or looking for focused learning, these tutorials get you up to speed fast.",
		"#tutorials",
		"Start Learning",
		user,
	)
	allTutorials := tutorials.GetAllTutorials()

	return &TutorialsPageData{
		HeaderComponentData: headerComponentData,
		AllTutorials:        allTutorials,
	}
}
