package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UsersAuth is for users authentication using email/password and generates a JWT token
// @Summary Login
// @Description Allows users to login to the platform and obtain a JWT Token
// @Tags Login
// @Accept  json
// @Produce  json
// @Param Account body models.Login true "account"
// @Success 200 {string} string
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /login [post]
func UsersAuth(c buffalo.Context) error {
	userauth := &models.UserAuth{User: models.User{}}
	//user := &models.User{}
	//cred := &models.Credential{}

	// Bind the credential to the JSON payload
	if err := c.Bind(userauth); err != nil {
		return errors.WithStack(err)
	}

	// Helper function to handle bad attempts
	bad := func() error {
		verrs := validate.NewErrors()
		return apiError(c, "Invalid email or password", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Fetch the user from the DB with the email
	tx := c.Value("tx").(*pop.Connection)
	err := loadUser(tx, userauth)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// Couldn't find an user with the supplied email address.
			return bad()
		}
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Confirm that the password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(userauth.Password), []byte(userauth.PasswordProvided))
	if err != nil {
		return bad()
	}

	// Create claims
	claims := jwt.MapClaims{}
	claims["Id"] = userauth.ID.String()
	claims["exp"] = time.Now().Add(oneWeek()).Unix()
	claims["role"] = userauth.User.Role
	claims["role_id"] = roleID(userauth)

	// Generate token
	secretKey := envy.Get("JWT_SECRET", "secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(fmt.Sprintf("{ \"token\": \"%s\" }", tokenString))))
}

func oneWeek() time.Duration {
	return 7 * 24 * time.Hour
}

// Set the role_id depending on the user role
func roleID(userauth *models.UserAuth) string {
	switch userauth.Role {
	case "admin":
		//claims["role_id"] = userauth.User.Admin.ID.String()
		return userauth.ID.String()
	case "parent":
		return userauth.User.Parent.ID.String()
	case "teacher":
		return userauth.User.Teacher.ID.String()
	default:
		return userauth.ID.String()

	}
}

func loadUser(tx *pop.Connection, userauth *models.UserAuth) error {
	var err error

	// Load the User first
	if err = tx.Where("email = ?", strings.ToLower(userauth.Email)).
		First(&userauth.User); err != nil {
		return err
	}

	// Load nested associations based on user Role
	switch userauth.User.Role {
	case "parent":
		parent := &models.Parent{}
		err = tx.Where("user_id = ?", userauth.User.ID).
			First(parent)
		userauth.User.Parent = parent
	case "teacher":
		teacher := &models.Teacher{}
		err = tx.Where("user_id = ?", userauth.User.ID).
			First(teacher)
		userauth.User.Teacher = teacher
	default: // default is admin
		//admin := &models.Admin{}
		//err = tx.Where("user_id = ?", userauth.User.ID).
		//First(admin)
		//userauth.User.Admin = admin
	}
	return err
}
