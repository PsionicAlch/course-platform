package config

import (
	"github.com/PsionicAlch/psionicalch-home/pkg/envloader"
	"github.com/PsionicAlch/psionicalch-home/pkg/envloader/validators"
)

const (
	development = "development"
	testing     = "testing"
	production  = "production"
)

func SetupConfig() error {
	variables := map[string]validators.ValidationFunc{
		"PORT":                       validators.NotEmptyValidator,
		"ENVIRONMENT":                validators.ChainValidators(validators.NotEmptyValidator, validators.InSliceValidator([]string{development, testing, production})),
		"DOMAIN_NAME":                validators.NotEmptyValidator,
		"SESSION_COOKIE_NAME":        validators.NotEmptyValidator,
		"AUTH_COOKIE_NAME":           validators.NotEmptyValidator,
		"AUTH_TOKEN_LIFETIME":        validators.ChainValidators(validators.NotEmptyValidator, validators.IntValidator),
		"EMAIL_TOKEN_LIFETIME":       validators.ChainValidators(validators.NotEmptyValidator, validators.IntValidator),
		"CURRENT_SECURE_COOKIE_KEY":  validators.NotEmptyValidator,
		"PREVIOUS_SECURE_COOKIE_KEY": validators.EmptyValidator,
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
