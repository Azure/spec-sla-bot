package actions

import (
	"context"
	"log"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/Azure/buffalo-azure/sdk/eventgrid"
	"github.com/Azure/spec-sla-bot/messages"
	"github.com/Azure/spec-sla-bot/models"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator
var commitID string

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_spec_sla_bot_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		//app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)

		//Subscribe to eventgrid
		eventgrid.RegisterSubscriber(app, "/specgithub", NewSpecgithubSubscriber(&eventgrid.BaseSubscriber{}))
		log.Printf("Created handler")
		//Create AMQP Listener
		messages.ReceiveFromQueue(context.Background(), os.Getenv("CUSTOMCONNSTR_SERVICEBUS_CONNECTION_STRING"))

		app.Resource("/assignees", AssigneesResource{})
		app.Resource("/events", EventsResource{})
		app.Resource("/emails", EmailsResource{})
		app.Resource("/email_assignees", EmailAssigneesResource{})
		app.Resource("/pullrequests", PullrequestsResource{})
		app.Resource("/pullrequest_assignees", PullrequestAssigneesResource{})
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}
