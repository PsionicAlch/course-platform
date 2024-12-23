package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewChangeFirstNameForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		FirstName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
	})
}

func ChangeFirstNamePartialValidation(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		FirstName: validators.ChainValidators(
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
	})
}

func EmptyChangeFirstNameFormComponent() *html.ChangeFirstNameFormComponent {
	firstNameInput := new(html.FormControlComponent)
	firstNameInput.Label = "First Name:"
	firstNameInput.Name = FirstName
	firstNameInput.Type = "text"
	firstNameInput.ValidationURL = ""

	changeFirstNameForm := new(html.ChangeFirstNameFormComponent)
	changeFirstNameForm.FirstNameInput = firstNameInput

	return changeFirstNameForm
}

func NewChangeFirstNameFormComponent(form *GenericForm) *html.ChangeFirstNameFormComponent {
	changeFirstNameFormComponent := EmptyChangeFirstNameFormComponent()

	changeFirstNameFormComponent.FirstNameInput.Value = form.GetValue(FirstName)
	changeFirstNameFormComponent.FirstNameInput.Errors = form.GetErrors(FirstName)

	return changeFirstNameFormComponent
}

func GetChangeFirstNameFormValues(form *GenericForm) (firstName string) {
	firstName = form.GetValue(FirstName)

	return
}
