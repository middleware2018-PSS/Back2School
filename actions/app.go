package actions

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/casbin/casbin"
	"github.com/cippaciong/jsonapi"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	popmw "github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/pop"
	tokenauth "github.com/middleware2018-PSS/back2_school/middleware/authentication"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/x/sessions"
	authorization "github.com/middleware2018-PSS/back2_school/middleware/authorization"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/rs/cors"

	// Use a blank import as required by swaggo
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

// App is the main of the API and here we define the General API info for swagger
// @title Back2School API
// @version 1.0
// @description This is an api to manage online services for a school.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/middleware2018-PSS/Back2School
// @contact.email sardelli.tommaso@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host localhost:3000
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
		model := envy.Get("AUTH_MODEL", "./auth_model.conf")
		policy := envy.Get("POLICY", "./policy.csv")
		authEnforcer, err := casbin.NewEnforcerSafe(model, policy, false)
		if err != nil {
			log.Fatal(err)
		}
		authorizer := authorization.New(authEnforcer)
		api.Use(authorizer)
		//api.Middleware.Skip(authorizer, UsersAuth)

		api.GET("/", listRoutes)
		api.POST("/login", UsersAuth)
		api.GET("/self", getSelf)
		api.Resource("/users", UsersResource{})
		api.GET("/users/{id}/{res:(?:notifications)}", func(c buffalo.Context) error {
			return getLists(c, &models.User{})
		})
		api.Resource("/admins", AdminsResource{})

		api.Resource("/parents", ParentsResource{})
		api.GET("/parents/{id}/{res:(?:students|appointments|payments|notifications)}", func(c buffalo.Context) error {
			return getLists(c, &models.Parent{})
		})
		api.POST("/parents/{id}/{res:(?:students)}", createNestedResource)

		api.Resource("/teachers", TeachersResource{})
		api.GET("/teachers/{id}/{res:(?:classes|appointments|notifications)}", func(c buffalo.Context) error {
			return getLists(c, &models.Teacher{})
		})
		api.POST("/teachers/{id}/{res:(?:appointments)}", createNestedResource)

		api.Resource("/students", StudentsResource{})
		api.GET("/students/{id}/{res:(?:parents|grades)}", func(c buffalo.Context) error {
			return getLists(c, &models.Student{})
		})

		api.Resource("/appointments", AppointmentsResource{})
		api.GET("/appointments/{id}/{res:(?:parents)}", func(c buffalo.Context) error {
			return getLists(c, &models.Appointment{})
		})

		api.Resource("/classes", ClassesResource{})
		api.GET("/classes/{id}/{res:(?:students|teachers|timetables)}", func(c buffalo.Context) error {
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

		api.Resource("/timetables", TimetablesResource{})

		api.GET("{all:.*}", func(c buffalo.Context) error {
			return apiError(c, "Route does not exist", "Not Found",
				http.StatusNotFound, errors.New("Route does not exist"))
		})
		api.POST("{all:.*}", func(c buffalo.Context) error {
			return apiError(c, "Route does not exist", "Not Found",
				http.StatusNotFound, errors.New("Route does not exist"))
		})
		api.PUT("{all:.*}", func(c buffalo.Context) error {
			return apiError(c, "Route does not exist", "Not Found",
				http.StatusNotFound, errors.New("Route does not exist"))
		})
		api.DELETE("{all:.*}", func(c buffalo.Context) error {
			return apiError(c, "Route does not exist", "Not Found",
				http.StatusNotFound, errors.New("Route does not exist"))
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

// List
func listRoutes(c buffalo.Context) error {
	var routesList []string
	for _, r := range app.Routes() {
		entry := fmt.Sprintf("%s %s", r.Method, r.Path)
		routesList = append(routesList, entry)
	}
	return c.Render(200, r.JSON(routesList))
}

func getSelf(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	// Get role and roleID from jwt claims
	var role, roleID string
	if claims, ok := c.Value("claims").(jwt.MapClaims); ok {
		role = claims["role"].(string)
		roleID = claims["role_id"].(string)
	} else {
		return apiError(c, "Token Error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("Token Error"))
	}

	res := new(bytes.Buffer)
	switch role {
	case "admin":
		baseres := &models.Admin{}
		if err := tx.Eager().Find(baseres, roleID); err != nil {
			return apiError(c, "Resource not found",
				"Not Found", http.StatusNotFound, err)
		}
		baseres.User.Password = ""
		err := jsonapi.MarshalPayload(res, baseres)
		if err != nil {
			return apiError(c, "Internal Error preparing the response payload",
				"Internal Server Error", http.StatusInternalServerError, err)
		}
	case "parent":
		baseres := &models.Parent{}
		if err := tx.Eager().Find(baseres, roleID); err != nil {
			return apiError(c, "Resource not found",
				"Not Found", http.StatusNotFound, err)
		}
		baseres.User.Password = ""
		err := jsonapi.MarshalPayload(res, baseres)
		if err != nil {
			return apiError(c, "Internal Error preparing the response payload",
				"Internal Server Error", http.StatusInternalServerError, err)
		}
	case "teacher":
		baseres := &models.Teacher{}
		if err := tx.Eager().Find(baseres, roleID); err != nil {
			return apiError(c, "Resource not found",
				"Not Found", http.StatusNotFound, err)
		}
		baseres.User.Password = ""
		err := jsonapi.MarshalPayload(res, baseres)
		if err != nil {
			return apiError(c, "Internal Error preparing the response payload",
				"Internal Server Error", http.StatusInternalServerError, err)
		}
	default:
		return apiError(c, "Token Error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("Token Error"))
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

func getLists(c buffalo.Context, baseres interface{}) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	resname := strings.Title(c.Param("res"))

	log.Println("RESNAME", resname)

	id := c.Param("id")

	if resname == "Notifications" {
		if err := tx.Find(baseres, id); err != nil {
			return apiError(c, "Cannot delete resource. Resource not found",
				"Not Found", http.StatusNotFound, err)
		}
		switch baseres.(type) {
		case *models.Parent:
			id = baseres.(*models.Parent).UserID.String()
		case *models.Teacher:
			id = baseres.(*models.Teacher).UserID.String()
		default:
			log.Println("SUCKA")
		}

		log.Println("ID SET:", id)
		baseres = &models.User{}
	}

	// To find the Parent the parameter parent_id is used.
	if err := tx.Eager(resname).Find(baseres, id); err != nil {
		return apiError(c, "Resource not found",
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
	case []*models.Grade:
		val = iface.([]*models.Grade)
	case []*models.Timetable:
		val = iface.([]*models.Timetable)
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
