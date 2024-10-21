package config

import (
	"time"

	"github.com/PsionicAlch/psionicalch-home/internal/utils"
	"github.com/joho/godotenv"
)

type ProjectConfig struct {
	port                         string
	environment                  string
	sessionCookieName            string
	domainName                   string
	authTokenLifetime            time.Duration
	emailTokenLifetime           time.Duration
	secureCookieHashKey          string
	secureCookieBlockKey         string
	secureCookiePreviousHashKey  string
	secureCookiePreviousBlockKey string
}

var projectConfig *ProjectConfig

const (
	development = "development"
	testing     = "testing"
	production  = "production"
)

// NewConfig creates a new singleton instance of ProjectConfig.
func NewConfig() {
	configLoggers := utils.CreateLoggers("CONFIG")

	// Load environment variables.
	err := godotenv.Load()
	if err != nil {
		configLoggers.ErrorLog.Fatalln("Failed to read .env file:", err)
	}

	// Get server port.
	port, err := GetVariable[string]("PORT")
	if err != nil {
		configLoggers.ErrorLog.Fatalln(err)
	}

	// Get project environment.
	environment, err := GetVariable[string]("ENVIRONMENT")
	if err != nil {
		configLoggers.ErrorLog.Fatalln(err)
	}

	viableEnvironments := []string{production, testing, development}
	if !utils.InSlice(environment, viableEnvironments) {
		configLoggers.ErrorLog.Fatalf("ENVIRONMENT (%s) can only be set to on of the these: %v!", environment, viableEnvironments)
	}

	// Get session cookie name.
	sessionCookieName, err := GetVariable[string]("SESSION_COOKIE_NAME")
	if err != nil {
		configLoggers.ErrorLog.Fatalln(err)
	}

	// Get domain name.
	domainName, err := GetVariable[string]("DOMAIN_NAME")
	if err != nil {
		configLoggers.ErrorLog.Fatalln(err)
	}

	// Get auth token lifetime.
	authTokenLifetime, err := GetVariable[time.Duration]("AUTH_TOKEN_LIFETIME")
	if err != nil {
		configLoggers.ErrorLog.Fatalln(err)
	}

	// Get email token lifetime.
	emailTokenLifetime, err := GetVariable[time.Duration]("EMAIL_TOKEN_LIFETIME")
	if err != nil {
		configLoggers.ErrorLog.Fatalln(err)
	}

	// Get securecookie hash key.
	secureCookieHashKey, err := GetVariable[string]("SECURE_COOKIE_HASH_KEY")
	if err != nil {
		configLoggers.ErrorLog.Fatalln(err)
	}

	projectConfig = &ProjectConfig{
		port:                port,
		environment:         environment,
		sessionCookieName:   sessionCookieName,
		domainName:          domainName,
		authTokenLifetime:   authTokenLifetime * time.Minute,
		emailTokenLifetime:  emailTokenLifetime * time.Minute,
		secureCookieHashKey: secureCookieHashKey,
	}
}

// GetPort returns the port.
func GetPort() string {
	return projectConfig.port
}

// GetEnvironment returns the environment.
func GetEnvironment() string {
	return projectConfig.environment
}

// InDevelopment states whether application environment is set to development.
func InDevelopment() bool {
	return projectConfig.environment == development
}

// InTesting states whether application environment is set to testing.
func InTesting() bool {
	return projectConfig.environment == testing
}

// InProduction states whether application environment is set to production.
func InProduction() bool {
	return projectConfig.environment == production
}

// GetSessionCookieName returns the name of the session cookie.
func GetSessionCookieName() string {
	return projectConfig.sessionCookieName
}

// GetDomainName returns the domain name of the application.
func GetDomainName() string {
	return projectConfig.domainName
}

// GetAuthTokenLifetime returns the authentication token's lifetime duration in minutes.
func GetAuthTokenLifetime() time.Duration {
	return projectConfig.authTokenLifetime
}

// GetEmailTokenLifetime returns the email token's lifetime duration in minutes.
func GetEmailTokenLifetime() time.Duration {
	return projectConfig.emailTokenLifetime
}
