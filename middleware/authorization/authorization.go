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
	switch role {
	case "admin":
		// This should not be necessary as we handle it already in auth_model.conf
		return true
	case "parent":
		return checkParentOwnership(roleID, rURL, pURL, role, c)
	case "teacher":
		return checkTeacherOwnership(roleID, rURL, pURL, role, c)
	}
	//if strings.Contains(pURL, ":id") &&
	//strings.Split(rURL, "/")[4] == roleID {
	//return true
	//}
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

func checkTeacherOwnership(roleID, rURL, pURL, role string, c *buffalo.DefaultContext) bool {
	log.Println("Checking teacher ownership")

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		log.Println("Database error in Casbin")
		return false
	}

	teacher := &models.Teacher{}

	basepath := strings.Split(rURL, "/")[3]
	switch basepath {
	case "classes":
		if err := tx.Eager("Classes").Find(teacher, roleID); err != nil {
			log.Println("Error loading teacher resource")
			return false
		}
		for _, class := range teacher.Classes {
			if class.ID.String() == c.Param("id") {
				return true
			}
		}
		return false
	default:
		return false
	}

}

func checkParentOwnership(roleID, rURL, pURL, role string, c *buffalo.DefaultContext) bool {
	return false
}
