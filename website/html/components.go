package html

import (
	"github.com/PsionicAlch/psionicalch-home/internal/database/models"
	"github.com/PsionicAlch/psionicalch-home/website/forms"
)

type HeaderComponentData struct {
	User     *models.UserModel
	Heading  string
	Text     string
	LinkHref string
	LinkText string
}

func CreateHeaderComponent(heading, text, linkhref, linktext string, user *models.UserModel) *HeaderComponentData {
	return &HeaderComponentData{
		User:     user,
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
