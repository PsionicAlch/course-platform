package config

import (
	"github.com/PsionicAlch/course-platform/pkg/envloader"
	"github.com/PsionicAlch/course-platform/pkg/envloader/validators"
)

const (
	development = "development"
	testing     = "testing"
	production  = "production"
)

func SetupConfig() error {
	variables := map[string]validators.ValidationFunc{
		"PORT":                     validators.NotEmpty,
		"ENVIRONMENT":              validators.InSlice([]string{development, testing, production}),
		"DOMAIN_NAME":              validators.NotEmpty,
		"NOTIFICATION_COOKIE_NAME": validators.NotEmpty,
		"AUTH_COOKIE_NAME":         validators.NotEmpty,
		"AUTH_TOKEN_LIFETIME": validators.Chain(
			validators.NotEmpty,
			validators.Int,
		),
		"EMAIL_TOKEN_LIFETIME": validators.Chain(
			validators.NotEmpty,
			validators.Int,
		),
		"CURRENT_SECURE_COOKIE_KEY":  validators.NotEmpty,
		"PREVIOUS_SECURE_COOKIE_KEY": validators.Empty,
		"EMAIL_PROVIDER":             validators.InSlice([]string{"smtp"}),
		"EMAIL_HOST":                 validators.NotEmpty,
		"EMAIL_PORT":                 validators.NotEmpty,
		"EMAIL_ADDRESS":              validators.NotEmpty,
		"EMAIL_PASSWORD":             validators.Empty,
		"STRIPE_SECRET_KEY":          validators.NotEmpty,
		"STRIPE_WEBHOOK_SECRET":      validators.NotEmpty,
		"CLOUDFRONT_URL":             validators.NotEmpty,
	}

	return envloader.LoadEnvironment(variables)
}

func Get[T any](name string) (T, error) {
	return envloader.GetVariable[T](name)
}

func GetWithoutError[T any](name string) T {
	variable, _ := envloader.GetVariable[T](name)
	return variable
}

func InDevelopment() bool {
	return GetWithoutError[string]("ENVIRONMENT") == development
}

func InTesting() bool {
	return GetWithoutError[string]("ENVIRONMENT") == testing
}

func InProduction() bool {
	return GetWithoutError[string]("ENVIRONMENT") == production
}
