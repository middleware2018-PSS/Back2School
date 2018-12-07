package authorization

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
)

func isOwner(roleID, rURL, pURL string) bool {
	if strings.Contains(pURL, ":id") &&
		strings.Split(rURL, "/")[4] == roleID {
		return true
	}
	return false
}

func isOwnerFunc(args ...interface{}) (interface{}, error) {
	roleID := args[0].(string)
	rURL := args[1].(string)
	pURL := args[2].(string)

	return (bool)(isOwner(roleID, rURL, pURL)), nil
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
			}

			e.AddFunction("isOwner", isOwnerFunc)

			// Casbin rule enforcing
			res, err := e.EnforceSafe(roleID, c.Value("current_path"), c.Value("method"), role)
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
