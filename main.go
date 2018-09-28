package main

import (
	"log"

	"github.com/middleware2018-PSS/back2_school/actions"
)

// main is the starting point to your Buffalo application.
// you can feel free and add to this `main` method, change
// what it does, etc...
// All we ask is that, at some point, you make sure to
// call `app.Serve()`, unless you don't want to start your
// application that is. :)

// @title Back2School API
// @version 1.0
// @description These are a School management system's API .
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:3000
func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}

/*
# Notes about `main.go`

## SSL Support

We recommend placing your application behind a proxy, such as
Apache or Nginx and letting them do the SSL heaving lifting
for you. https://gobuffalo.io/en/docs/proxy

## Buffalo Build

When `buffalo build` is run to compile your binary this `main`
function will be at the heart of that binary. It is expected
that your `main` function will start your application using
the `app.Serve()` method.

*/
