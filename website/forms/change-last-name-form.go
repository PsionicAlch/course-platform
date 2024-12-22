package forms

import (
	"net/http"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
	"github.com/PsionicAlch/psionicalch-home/website/html"
)

func NewChangeLastNameForm(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		LastName: validators.ChainValidators(
			validators.NotEmpty,
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
	})
}

func ChangeLastNamePartialValidation(r *http.Request) *GenericForm {
	return NewForm(r, map[FieldName]validators.ValidationFunc{
		LastName: validators.ChainValidators(
			validators.MaxLength(50),
			validators.Profanities,
			validators.NoSpecialCharacters,
			validators.NoNumbers,
		),
	})
}

func EmptyChangeLastNameFormComponent() *html.ChangeLastNameFormComponent {
	lastNameInput := new(html.FormControlComponent)
	lastNameInput.Label = "Last Name:"
	lastNameInput.Name = LastName
	lastNameInput.Type = "text"
	lastNameInput.ValidationURL = ChangeLastNameValidationURL

	changeLastNameForm := new(html.ChangeLastNameFormComponent)
	changeLastNameForm.LastNameInput = lastNameInput

	return changeLastNameForm
}

func NewChangeLastNameFormComponent(form *GenericForm) *html.ChangeLastNameFormComponent {
	changeLastNameFormComponent := EmptyChangeLastNameFormComponent()

	changeLastNameFormComponent.LastNameInput.Value = form.GetValue(LastName)
	changeLastNameFormComponent.LastNameInput.Errors = form.GetErrors(LastName)

	return changeLastNameFormComponent
}

func GetChangeLastNameFormValues(form *GenericForm) (lastName string) {
	lastName = form.GetValue(LastName)

	return
}
