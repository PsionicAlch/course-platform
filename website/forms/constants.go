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
	AffiliateCodeName   = "affiliate_code"
	AffiliatePointsName = "affiliate_points"
	DiscountCodeName    = "discount_code"

	// Validation URLs
	SignupValidationURL        = "/accounts/validate/signup"
	ResetPasswordValidationURL = "/accounts/validate/reset-password"
	NewDiscountsValidationURL  = "/admin/discounts/validate/add"
)
