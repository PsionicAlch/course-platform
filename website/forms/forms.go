package forms

import (
	"net/http"
	"net/url"

	"github.com/PsionicAlch/psionicalch-home/website/forms/validators"
)

type FieldName string

type GenericForm struct {
	values      url.Values
	errors      map[FieldName][]string
	validations map[FieldName]validators.ValidationFunc
}

func NewForm(r *http.Request, validations map[FieldName]validators.ValidationFunc) *GenericForm {
	r.ParseForm()

	return &GenericForm{
		values:      r.Form,
		errors:      make(map[FieldName][]string),
		validations: validations,
	}
}

func (form *GenericForm) Validate() {
	form.errors = make(map[FieldName][]string)

	for fieldName, validatorFunc := range form.validations {
		err := validatorFunc(form.values.Get(string(fieldName)), form.values)
		if err != nil {
			form.errors[fieldName] = append(form.errors[fieldName], err.Error())
		}
	}
}

func (form *GenericForm) GetValue(field FieldName) string {
	return form.values.Get(string(field))
}

func (form *GenericForm) GetErrors(field FieldName) []string {
	errors, contains := form.errors[field]
	if contains {
		return errors
	}

	return []string{}
}
