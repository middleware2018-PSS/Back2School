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
)

type Admin struct {
	ID        uuid.UUID `json:"id" db:"id" jsonapi:"primary,admins"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	Email     string    `json:"email" db:"email" jsonapi:"attr,email"`
	Password  string    `json:"passowrd" db:"-" jsonapi:"attr,password,omitempty"`
	Name      string    `json:"name" db:"name" jsonapi:"attr,name"`
	Surname   string    `json:"surname" db:"surname" jsonapi:"attr,surname"`
	UserID    uuid.UUID `db:"user_id"`
	User      *User     `db:"-" jsonapi:"relation,user,omitempty"`
}

// String is not required by pop and may be deleted
func (a Admin) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Admins is not required by pop and may be deleted
type Admins []Admin

// String is not required by pop and may be deleted
func (a Admins) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Admin) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Field: a.Email, Name: "Email", Message: mailValidationMsg},
		&validators.StringIsPresent{Field: a.Name, Name: "Name", Message: nameValidationMsg},
		&validators.StringIsPresent{Field: a.Surname, Name: "Surname", Message: surnameValidationMsg},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Admin) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Admin) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a admin
func (admin Admin) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/admins/%s", APIUrl, admin.ID.String()),
	}
}

// Invoked for each relationship defined on the Admin struct when marshaled
func (admin Admin) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "user" {
		return &jsonapi.Links{
			"user": fmt.Sprintf("http://%s/users/%s", APIUrl, admin.UserID.String()),
		}
	}
	return nil
}
