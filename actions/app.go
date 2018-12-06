package actions

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/casbin/casbin"
	"github.com/cippaciong/jsonapi"
	"github.com/gobuffalo/buffalo"
	popmw "github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	tokenauth "github.com/gobuffalo/mw-tokenauth"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/x/sessions"
	authorization "github.com/middleware2018-PSS/back2_school/middleware/authorization"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/rs/cors"

	_ "github.com/middleware2018-PSS/back2_school/docs"
	buffaloSwagger "github.com/swaggo/buffalo-swagger"
	"github.com/swaggo/buffalo-swagger/swaggerFiles"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// Create a new instance of the logger
var log = logrus.New()

// @title Back2School API
// @version 1.0
// @description This is an api to manage online services for a school.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/middleware2018-PSS/Back2School
// @contact.email sardelli.tommaso@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /api/v1
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

		//app.GET("/", HomeHandler)
		app.Redirect(301, "/", "/swagger/index.html")
		app.GET("/swagger/{doc:.*}", buffaloSwagger.WrapHandler(swaggerFiles.Handler))

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

		api.GET("/", ListRoutes)
		api.POST("/login", UsersAuth)
		api.Resource("/users", UsersResource{})
		api.GET("/users/{id}/{res:(?:notifications)}", func(c buffalo.Context) error {
			return getLists(c, &models.User{})
		})
		api.Resource("/admins", AdminsResource{})

		api.Resource("/parents", ParentsResource{})
		api.GET("/parents/{id}/{res:(?:students|appointments|payments)}", func(c buffalo.Context) error {
			return getLists(c, &models.Parent{})
		})
		api.POST("/parents/{id}/{res:(?:students)}", createNestedResource)

		api.Resource("/teachers", TeachersResource{})
		api.GET("/teachers/{id}/{res:(?:classes|appointments)}", func(c buffalo.Context) error {
			return getLists(c, &models.Teacher{})
		})
		api.POST("/teachers/{id}/{res:(?:appointments)}", createNestedResource)

		api.Resource("/students", StudentsResource{})
		api.GET("/students/{id}/{res:(?:parents)}", func(c buffalo.Context) error {
			return getLists(c, &models.Student{})
		})

		api.Resource("/appointments", AppointmentsResource{})
		api.GET("/appointments/{id}/{res:(?:parents)}", func(c buffalo.Context) error {
			return getLists(c, &models.Appointment{})
		})

		api.Resource("/classes", ClassesResource{})
		api.GET("/classes/{id}/{res:(?:students|teachers)}", func(c buffalo.Context) error {
			return getLists(c, &models.Class{})
		})

		api.Resource("/grades", GradesResource{})

		api.Resource("/notifications", NotificationsResource{})
		api.GET("/notifications/{id}/{res:(?:users)}", func(c buffalo.Context) error {
			return getLists(c, &models.Notification{})
		})

		api.Resource("/payments", PaymentsResource{})
		api.GET("/payments/{id}/{res:(?:parents)}", func(c buffalo.Context) error {
			return getLists(c, &models.Payment{})
		})

		api.POST("/payments/{payment_id}/pay", FakePay)
		api.GET("{all:.*}", func(c buffalo.Context) error {
			return c.Render(200, r.Func("application/json",
				customJSONRenderer("404 Not found")))
		})
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

func ListRoutes(c buffalo.Context) error {
	var routesList []string
	for _, r := range app.Routes() {
		entry := fmt.Sprintf("%s %s", r.Method, r.Path)
		routesList = append(routesList, entry)
	}
	return c.Render(200, r.JSON(routesList))
}

func getLists(c buffalo.Context, baseres interface{}) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	resname := strings.Title(c.Param("res"))

	// To find the Parent the parameter parent_id is used.
	if err := tx.Eager(resname).Find(baseres, c.Param("id")); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	iface := reflect.ValueOf(baseres).Elem().FieldByName(resname).Interface()
	var val interface{}
	switch iface.(type) {
	case []*models.User:
		val = iface.([]*models.User)
		for _, u := range val.([]*models.User) {
			(*u).Password = ""
		}
	case []*models.Parent:
		val = iface.([]*models.Parent)
	case []*models.Student:
		val = iface.([]*models.Student)
	case []*models.Teacher:
		val = iface.([]*models.Teacher)
	case []*models.Appointment:
		val = iface.([]*models.Appointment)
	case []*models.Notification:
		val = iface.([]*models.Notification)
	case []*models.Payment:
		val = iface.([]*models.Payment)
	case []*models.Class:
		val = iface.([]*models.Class)
	default:
		log.Println("SUCKA")
	}

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, val)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

func createNestedResource(c buffalo.Context) error {
	var val interface{}
	switch c.Param("res") {
	case "students":
		val = &models.Student{}
	case "appointments":
		val = &models.Appointment{}
	default:
		log.Println("SUCKA")
	}

	if err := jsonapi.UnmarshalPayload(c.Request().Body, val); err != nil {
		log.Println("ERROR UNMARSHALLING")
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Create and save the student
	verrs, err := tx.ValidateAndCreate(val)
	if err != nil {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, err)
	}

	// Check for validation errors
	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Reload the resource with all the relationships

	if c.Param("res") == "students" {
		if err := tx.Eager("Parents").Find(val, val.(*models.Student).ID); err != nil {
			return apiError(c, "The requested resource cannot be found",
				"Not Found", http.StatusNotFound, err)
		}
	} else {
		if err := tx.Eager().Find(val, val.(*models.Appointment).ID); err != nil {
			return apiError(c, "The requested resource cannot be found",
				"Not Found", http.StatusNotFound, err)
		}
	}

	// If there are no errors return the Student resource
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, val)
	if err != nil {
		log.Println("ERROR MARSHALLING")
		return apiError(c, "Error processing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}
	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}
