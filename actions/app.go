package actions

import (
	"github.com/gobuffalo/buffalo"
	popmw "github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	tokenauth "github.com/gobuffalo/mw-tokenauth"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/x/sessions"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

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
		app.Use(forceSSL())

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
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

		api.POST("/login", UsersAuth)
		api.Resource("/users", UsersResource{})
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
