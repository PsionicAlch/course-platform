package forms

import "github.com/PsionicAlch/psionicalch-home/internal/validators"

type Form interface {
	GetErrors() map[string][]string
	SetErrors(errs map[string][]string)
}

func AppendErrors(err error, field string, form Form) {
	formErrors := form.GetErrors()

	switch e := err.(type) {
	case *validators.ValidationError:
		formErrors[field] = append(formErrors[field], e.Errors...)
	default:
		formErrors[field] = append(formErrors[field], "an unknown error has occurred")
	}

	form.SetErrors(formErrors)
}

func AppendError(field, message string, form Form) {
	formErrors := form.GetErrors()
	formErrors[field] = append(formErrors[field], message)

	form.SetErrors(formErrors)
}
