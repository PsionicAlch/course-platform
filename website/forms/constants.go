package forms

const (
	// Form Field Names
	FirstName           = "first_name"
	LastName            = "last_name"
	EmailName           = "email"
	PasswordName        = "password"
	ConfirmPasswordName = "confirm_password"
	TitleName           = "title"
	DescriptionName     = "description"
	UsesName            = "uses"
	AmountName          = "amount"

	// Validation URLs
	SignupValidationURL        = "/accounts/validate/signup"
	ResetPasswordValidationURL = "/accounts/validate/reset-password"
	NewDiscountsValidationURL  = "/admin/discounts/validate/add"
)
