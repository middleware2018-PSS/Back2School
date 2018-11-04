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

func isOwner(role_id, r_url, p_url string) bool {
	if strings.Contains(p_url, ":id") &&
		strings.Split(r_url, "/")[4] == role_id {
		return true
	}
	return false
}

func isOwnerFunc(args ...interface{}) (interface{}, error) {
	role_id := args[0].(string)
	r_url := args[1].(string)
	p_url := args[2].(string)

	return (bool)(isOwner(role_id, r_url, p_url)), nil
}

func New(e *casbin.Enforcer) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		fn := func(c buffalo.Context) error {
			var role, role_id string
			if claims, ok := c.Value("claims").(jwt.MapClaims); ok {
				role = claims["role"].(string)
				role_id = claims["role_id"].(string)
			} else {
				role = "anonymous"
			}

			e.AddFunction("isOwner", isOwnerFunc)

			// Casbin rule enforcing
			res, err := e.EnforceSafe(role_id, c.Value("current_path"), c.Value("method"), role)
			if err != nil {
				log.Println("Error loading Casbin enforcing")
				log.Println(err)
				return c.Error(http.StatusInternalServerError, err)
			}
			if res {
				log.Println("It's ok to go on with buffalo")
				err = next(c)
			} else {
				return c.Error(http.StatusForbidden, errors.New("You are not authorized to do this"))
			}

			log.Println("Ready to return buffalo err")
			return err
		}

		return fn
	}
}
