package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cippaciong/jsonapi"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserAuth is a wrapper around User used for authentication
type UserAuth struct {
	User
	Email            string `json:"email" db:"-"`
	PasswordProvided string `json:"password" db:"-"`
}

// User is the model for users registered to the school system (admins, parents or teachers)
type User struct {
	ID            uuid.UUID       `json:"id" db:"id" jsonapi:"primary,users"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	Email         string          `json:"email" db:"email" jsonapi:"attr,email"`
	Password      string          `json:"password" db:"password" jsonapi:"attr,password,omitempty"`
	Role          string          `json:"role" db:"role" jsonapi:"attr,role"`
	Parent        *Parent         `has_one:"parent" fk_id:"user_id" jsonapi:"relation,parent,omitempty"`
	Teacher       *Teacher        `has_one:"teacher" fk_id:"user_id" jsonapi:"relation,teacher,omitempty"`
	Admin         *Admin          `has_one:"admin" fk_id:"user_id" jsonapi:"relation,admin,omitempty"`
	Notifications []*Notification `many_to_many:"users_notifications" jsonapi:"relation,notifications,omitempty"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []*User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// CALLBACKS //

// BeforeCreate encrypts user password with bcrypt before
// storing it in the database
func (u *User) BeforeCreate(tx *pop.Connection) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	u.Password = string(hash)

	return nil
}

// VALIDATIONS //

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Field: u.Email, Name: "Email", Message: mailValidationMsg},
		&validators.StringIsPresent{Field: u.Password, Name: "Password", Message: passwordValidationMsg},
		&validators.StringIsPresent{Field: u.Role, Name: "Role", Message: roleValidationMsg},
	), nil
}

// JSONAPILinks implements the Linkable interface for a user
func (u User) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/users/%s", APIUrl, u.ID.String()),
	}
}

// JSONAPIRelationshipLinks is invoked for each relationship defined on the User struct when marshaled
func (u User) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "notifications" {
		return &jsonapi.Links{
			"notifications": fmt.Sprintf("http://%s/users/%s/notifications",
				APIUrl, u.ID.String()),
		}
	}
	if relation == "admin" {
		return &jsonapi.Links{
			"admin": fmt.Sprintf("http://%s/admins/%s",
				APIUrl, u.Admin.ID.String()),
		}
	}
	return nil
}
