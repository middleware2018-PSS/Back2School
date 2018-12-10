package authorization

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
)

func isOwner(roleID, rURL, pURL, role string, c *buffalo.DefaultContext) bool {
	// Get the name of the resource we are accessing
	res := strings.Split(rURL, "/")[3]

	// Direct check if parent is accessing a parent resource
	if res == "parents" {
		return c.Param("parent_id") == roleID || c.Param("id") == roleID
	}

	// Direct check if teacher is accessing a teacher resource
	if res == "teachers" {
		return c.Param("teacher_id") == roleID || c.Param("id") == roleID
	}

	// For all other endpoints exploit the Ownable interface
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		log.Println("Database error in Casbin")
		return false
	}

	var id string
	var o models.Ownable
	switch res {
	case "students":
		if c.Param("student_id") != "" {
			id = c.Param("student_id")
		} else {
			id = c.Param("id")
		}
		o = &models.Student{}
		if err := tx.Eager().Find(o, id); err != nil {
			log.Println("Error eager loading")
			return false
		}
	case "payments":
		o = &models.Payment{}
		if err := tx.Eager().Find(o, c.Param("payment_id")); err != nil {
			log.Println("Error eager loading")
			return false
		}
	case "classes":
		if c.Param("class_id") != "" {
			id = c.Param("class_id")
		} else {
			id = c.Param("id")
		}
		o = &models.Class{}
		if err := tx.Eager().Find(o, id); err != nil {
			log.Println("Error eager loading")
			return false
		}
	case "grades":
		o = &models.Grade{}
		if err := tx.Eager().Find(o, c.Param("grade_id")); err != nil {
			log.Println("Error eager loading")
			return false
		}
	case "notifications":
		o = &models.Notification{}
		if err := tx.Eager().Find(o, c.Param("notification_id")); err != nil {
			log.Println("Error eager loading")
			return false
		}
	case "appointments":
		o = &models.Appointment{}
		if err := tx.Eager().Find(o, c.Param("appointment_id")); err != nil {
			log.Println("Error eager loading")
			return false
		}
	default:
		log.Println("SUCKA")
		return false
	}
	if role == "parent" {
		return o.BelongsToParent(tx, roleID)
	}
	if role == "teacher" {
		return o.BelongsToTeacher(tx, roleID)
	}

	return false
}

func isOwnerFunc(args ...interface{}) (interface{}, error) {
	roleID := args[0].(string)
	rURL := args[1].(string)
	pURL := args[2].(string)
	role := args[3].(string)
	c := args[4].(*buffalo.DefaultContext)

	return (bool)(isOwner(roleID, rURL, pURL, role, c)), nil
}

// New creates a new Buffalo Middleware for Casbin
func New(e *casbin.Enforcer) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		fn := func(c buffalo.Context) error {
			var role, roleID string
			if claims, ok := c.Value("claims").(jwt.MapClaims); ok {
				role = claims["role"].(string)
				roleID = claims["role_id"].(string)
			} else {
				role = "anonymous"
				roleID = ""
			}

			e.AddFunction("isOwner", isOwnerFunc)

			// Casbin rule enforcing
			log.Println("r.obj:", c.Value("current_path"))
			res, err := e.EnforceSafe(roleID, c.Value("current_path"), c.Value("method"), role, c)
			if err != nil {
				log.Println("Error loading Casbin enforcing")
				log.Println(err)
				return c.Error(http.StatusInternalServerError, err)
			}
			if res {
				err = next(c)
			} else {
				return c.Error(http.StatusForbidden, errors.New("You are not authorized to do this"))
			}
			return err
		}

		return fn
	}
}
