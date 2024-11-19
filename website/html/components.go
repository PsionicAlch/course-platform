package html

import "github.com/PsionicAlch/psionicalch-home/internal/database/models"

type NavbarComponent struct {
	User *models.UserModel
}

func NewNavbarComponent(user *models.UserModel) *NavbarComponent {
	return &NavbarComponent{
		User: user,
	}
}

/*
FormControlComponent is a data type to represent the properties required
to render a form-control component.

Form Control Props Required:

  - Label (a string for the label)

  - Type (a string to control the input type)

  - Name (a string to be used as the input name, id and label for)

  - Errors (a string array of errors)

  - ValidationURL (a string pointing to the URL that the input can use to validate)

  - Value (a string to be used as input value)
*/
type FormControlComponent struct {
	Label         string
	Type          string
	Name          string
	Errors        []string
	ValidationURL string
	Value         string
}

/*
PasswordControlComponent is a data type to represent the properties required
to render a form-control component.

Form Control Props Required:

  - Label (a string for the label)

  - Name (a string to be used as the input name, id and label for)

  - Errors (a string array of errors)

  - ValidationURL (a string pointing to the URL that the input can use to validate)

  - Value (a string to be used as input value)
*/
type PasswordControlComponent struct {
	Label         string
	Name          string
	Errors        []string
	ValidationURL string
	Value         string
}

type LoginFormComponent struct {
	EmailInput    *FormControlComponent
	PasswordInput *PasswordControlComponent
}

type SignupFormComponent struct {
	FirstNameInput       *FormControlComponent
	LastNameInput        *FormControlComponent
	EmailInput           *FormControlComponent
	PasswordInput        *PasswordControlComponent
	ConfirmPasswordInput *PasswordControlComponent
}

type ForgotPasswordFormComponent struct {
	EmailInput *FormControlComponent
}

type ResetPasswordFormComponent struct {
	EmailToken           string
	PasswordInput        *PasswordControlComponent
	ConfirmPasswordInput *PasswordControlComponent
}

type TutorialsListComponent struct {
	Tutorials    []*models.TutorialModel
	LastTutorial *models.TutorialModel
	QueryURL     string
	ErrorMessage string
}
