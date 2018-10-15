package authorization

import (
	"errors"
	"net/http"

	"github.com/casbin/casbin"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
)

func New(e *casbin.Enforcer) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		fn := func(c buffalo.Context) error {
			var role string
			if claims, ok := c.Value("claims").(jwt.MapClaims); ok {
				role = claims["role"].(string)
			} else {
				role = "anonymous"
			}

			// Casbin rule enforcing
			res, err := e.EnforceSafe(role, c.Value("current_path"), c.Value("method"))
			if err != nil {
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
