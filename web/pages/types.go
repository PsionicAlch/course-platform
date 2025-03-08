package pages

import (
	"fmt"
	"time"

	"github.com/PsionicAlch/course-platform/internal/authentication"
	"github.com/PsionicAlch/course-platform/internal/cache"
	gocache "github.com/PsionicAlch/course-platform/internal/cache/go-cache"
	"github.com/PsionicAlch/course-platform/internal/database"
	"github.com/PsionicAlch/course-platform/internal/database/sqlite_database"
	"github.com/PsionicAlch/course-platform/internal/payments"
	"github.com/PsionicAlch/course-platform/internal/render"
	vanillahtml "github.com/PsionicAlch/course-platform/internal/render/renderers/vanilla_html"
	vanillatext "github.com/PsionicAlch/course-platform/internal/render/renderers/vanilla_text"
	"github.com/PsionicAlch/course-platform/internal/session"
	"github.com/PsionicAlch/course-platform/internal/utils"
	"github.com/PsionicAlch/course-platform/web/config"
	"github.com/PsionicAlch/course-platform/web/emails"
	"github.com/PsionicAlch/course-platform/web/generators"
	"github.com/PsionicAlch/course-platform/web/html"
	"github.com/PsionicAlch/sitemapper"
)

type Renderers struct {
	Page render.Renderer
	Htmx render.Renderer
	RSS  render.Renderer
}

func CreateRenderers(pageRenderer render.Renderer, htmxRenderer render.Renderer, rssRenderer render.Renderer) *Renderers {
	return &Renderers{
		Page: pageRenderer,
		Htmx: htmxRenderer,
		RSS:  rssRenderer,
	}
}

type HandlerContext struct {
	Renderers      *Renderers
	Database       database.Database
	Authentication *authentication.Authentication
	Session        *session.Session
	Payment        *payments.Payments
	Emailer        *emails.Emails
	Cache          cache.Cache
	Mapper         *sitemapper.SiteMapper
}

func CreateHandlerContext() (*HandlerContext, error) {
	// Set up notifications.
	sessions := SetupSession()

	// Set up renderers.
	renderers, err := SetupRenderers(sessions)
	if err != nil {
		return nil, err
	}

	// Set up database.
	db, err := SetupDatabase()
	if err != nil {
		return nil, err
	}

	// Set up authentication system.
	auth, err := SetupAuthentication(db, sessions)
	if err != nil {
		return nil, err
	}

	// Set up emailer.
	emailer, err := SetupEmailer()
	if err != nil {
		return nil, err
	}

	// Set up payments.
	payment := SetupPayments(db, emailer)

	// Set up cache.
	cache := SetupCache(db, renderers.RSS)

	// Set up sitemapper.
	mapper := SetupSiteMapper()

	context := &HandlerContext{
		Renderers:      renderers,
		Database:       db,
		Authentication: auth,
		Session:        sessions,
		Payment:        payment,
		Emailer:        emailer,
		Cache:          cache,
		Mapper:         mapper,
	}

	return context, nil
}

func SetupSession() *session.Session {
	cookieName := config.GetWithoutError[string]("NOTIFICATION_COOKIE_NAME")
	domainName := config.GetWithoutError[string]("DOMAIN_NAME")
	return session.SetupSession(cookieName, domainName)
}

func SetupRenderers(sessions *session.Session) (*Renderers, error) {
	cloudfrontURL := config.GetWithoutError[string]("CLOUDFRONT_URL")
	pagesRenderer, err := vanillahtml.SetupVanillaHTMLRenderer(cloudfrontURL, sessions, html.HTMLFiles, ".page.tmpl", "pages", "layouts/*.layout.tmpl", "components/*.component.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to set up pages renderer: %w", err)
	}

	htmxRenderer, err := vanillahtml.SetupVanillaHTMLRenderer(cloudfrontURL, nil, html.HTMLFiles, ".htmx.tmpl", "htmx", "components/*.component.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to set up htmx renderer: %w", err)
	}

	xmlRenderer, err := vanillatext.SetupVanillaTextRenderer(cloudfrontURL, html.XMLFiles, ".xml.tmpl", "xml")
	if err != nil {
		return nil, fmt.Errorf("failed to set up rss renderer: %w", err)
	}

	return CreateRenderers(pagesRenderer, htmxRenderer, xmlRenderer), nil
}

func SetupDatabase() (database.Database, error) {
	db, err := sqlite_database.CreateSQLiteDatabase("/db/db.sqlite", "/db/migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	return db, nil
}

func SetupAuthentication(db database.Database, sessions *session.Session) (*authentication.Authentication, error) {
	authLifetime := time.Duration(config.GetWithoutError[int]("AUTH_TOKEN_LIFETIME")) * time.Minute
	pwdResetLifetime := time.Minute * 30
	domainName := config.GetWithoutError[string]("DOMAIN_NAME")
	cookieName := config.GetWithoutError[string]("AUTH_COOKIE_NAME")
	currentKey := config.GetWithoutError[string]("CURRENT_SECURE_COOKIE_KEY")
	previousKey := config.GetWithoutError[string]("PREVIOUS_SECURE_COOKIE_KEY")
	auth, err := authentication.SetupAuthentication(db, sessions, authLifetime, pwdResetLifetime, cookieName, domainName, currentKey, previousKey)
	if err != nil {
		return nil, fmt.Errorf("failed to set up authentication: %w", err)
	}

	return auth, nil
}

func SetupPayments(db database.Database, emailer *emails.Emails) *payments.Payments {
	stripeSecretKey := config.GetWithoutError[string]("STRIPE_SECRET_KEY")
	stripeWebhookSecret := config.GetWithoutError[string]("STRIPE_WEBHOOK_SECRET")

	return payments.SetupPayments(stripeSecretKey, stripeWebhookSecret, db, emailer)
}

func SetupEmailer() (*emails.Emails, error) {
	cloudfrontURL := config.GetWithoutError[string]("CLOUDFRONT_URL")
	emailRenderer, err := vanillahtml.SetupVanillaHTMLRenderer(cloudfrontURL, nil, html.HTMLFiles, ".email.tmpl", "emails", "layouts/email.layout.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to set up email renderer: %w", err)
	}

	return emails.SetupEmails(emailRenderer), nil
}

func SetupCache(db database.Database, xmlRenderer render.Renderer) cache.Cache {
	gens := &gocache.Generators{
		RSSFeed:                generators.RSSFeed(utils.CreateLoggers("RSS FEED GENERATOR"), db, xmlRenderer),
		TutorialsRSSFeed:       generators.TutorialsRSSFeed(utils.CreateLoggers("TUTORIALS RSS FEED GENERATOR"), db, xmlRenderer),
		TutorialRssFeed:        generators.TutorialRSSFeed(utils.CreateLoggers("TUTORIAL RSS FEED GENERATOR"), db, xmlRenderer),
		CoursesRSSFeed:         generators.CoursesRSSFeed(utils.CreateLoggers("COURSES RSS FEED GENERATOR"), db, xmlRenderer),
		AuthorTutorialsRSSFeed: generators.AuthorTutorialsRSSFeed(utils.CreateLoggers("AUTHOR TUTORIALS RSS FEED GENERATOR"), db, xmlRenderer),
		AuthorCoursesRSSFeed:   generators.AuthorCoursesRSSFeed(utils.CreateLoggers("AUTHOR COURSES RSS FEED GENERATOR"), db, xmlRenderer),
	}
	return gocache.SetupGoCache(gens)
}

func SetupSiteMapper() *sitemapper.SiteMapper {
	loggers := utils.CreateLoggers("SITEMAPPER")

	options := sitemapper.DefaultOptions()
	options.SetLinkAttributes("hx-get")
	options.SetErrorLogger(func(err error) {
		loggers.ErrorLog.Println(err)
	})

	return sitemapper.NewSiteMapper(options)
}
