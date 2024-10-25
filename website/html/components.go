package html

import "github.com/PsionicAlch/psionicalch-home/website/forms"

type HeaderComponentData struct {
	Heading  string
	Text     string
	LinkHref string
	LinkText string
}

func CreateHeaderComponent(heading, text, linkhref, linktext string) *HeaderComponentData {
	return &HeaderComponentData{
		Heading:  heading,
		Text:     text,
		LinkHref: linkhref,
		LinkText: linktext,
	}
}

type SignupFormComponentData struct {
	Form  *forms.SignUpForm
	Error string
}

type LoginFormComponentData struct {
	Form  *forms.LoginForm
	Error string
}
