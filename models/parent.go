package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cippaciong/jsonapi"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Parent struct {
	ID        uuid.UUID `json:"id" db:"id" jsonapi:"primary,parents"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	Email     string    `json:"email" db:"email" jsonapi:"attr,email"`
	Password  string    `json:"passowrd" db:"-" jsonapi:"attr,password,omitempty"`
	Name      string    `json:"name" db:"name" jsonapi:"attr,name"`
	Surname   string    `json:"surname" db:"surname" jsonapi:"attr,surname"`
	UserID    uuid.UUID `db:"user_id"` // jsonapi:"relation,user,omitempty"`
	User      *User     `db:"-" jsonapi:"relation,user,omitempty"`
}

// Return a string representation of the Parent resource
func (p Parent) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Parents is not required by pop and may be deleted
type Parents []Parent

// Return a string representation of the Parents resource
func (p Parents) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *Parent) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Field: p.Email, Name: "Email",
			Message: "Mail is not in the right format."},
		//&validators.StringIsPresent{Field: p.Password, Name: "Passowrd"},
		&validators.StringIsPresent{Field: p.Name, Name: "Name"},
		&validators.StringIsPresent{Field: p.Surname, Name: "Surname"},
	), nil
}

// JSONAPILinks implements the Linkable interface for a parent
func (parent Parent) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/parents/%s", APIUrl, parent.ID.String()),
	}
}

// Invoked for each relationship defined on the Parent struct when marshaled
func (parent Parent) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	log.Println()
	log.Println("Building Links")
	log.Println()
	if relation == "user" {
		log.Println(parent.UserID.String())
		return &jsonapi.Links{
			"user": fmt.Sprintf("http://%s/users/%s", APIUrl, parent.UserID.String()),
		}
	}
	return nil
}
