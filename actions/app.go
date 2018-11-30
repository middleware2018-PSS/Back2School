package actions

import (
	"github.com/casbin/casbin"
	"github.com/gobuffalo/buffalo"
	popmw "github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	tokenauth "github.com/gobuffalo/mw-tokenauth"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/x/sessions"
	authorization "github.com/middleware2018-PSS/back2_school/middleware/authorization"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// Create a new instance of the logger
var log = logrus.New()

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_back2_school_session",
		})
		// Automatically redirect to SSL
		//app.Use(forceSSL())

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
			log.SetLevel(logrus.TraceLevel)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		app.GET("/", HomeHandler)

		api := app.Group("/api/v1")

		// JWT Auth middleware
		TokenAuth := tokenauth.New(tokenauth.Options{})
		api.Use(TokenAuth)
		api.Middleware.Skip(TokenAuth, UsersAuth)

		// Setup casbin auth rules
		authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
		if err != nil {
			log.Fatal(err)
		}
		authorizer := authorization.New(authEnforcer)
		api.Use(authorizer)
		//api.Middleware.Skip(authorizer, UsersAuth)

		api.POST("/login", UsersAuth)
		api.Resource("/users", UsersResource{})
		api.Resource("/admins", AdminsResource{})
		api.Resource("/parents", ParentsResource{})
		api.Resource("/teachers", TeachersResource{})
		api.Resource("/students", StudentsResource{})
		//t.Resource("/students", StudentsResource{}) // parents->students nested resource
		//app.Resource("/students", StudentsResource{})
		api.Resource("/appointments", AppointmentsResource{})
		api.Resource("/classes", ClassesResource{})
		api.Resource("/grades", GradesResource{})
		api.Resource("/notifications", NotificationsResource{})
	}

	return app
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
