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

// Parent is the model for the parents of school students
type Parent struct {
	ID           uuid.UUID      `json:"id" db:"id" jsonapi:"primary,parents"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	Email        string         `json:"email" db:"email" jsonapi:"attr,email"`
	Password     string         `json:"passowrd" db:"-" jsonapi:"attr,password,omitempty"`
	Name         string         `json:"name" db:"name" jsonapi:"attr,name"`
	Surname      string         `json:"surname" db:"surname" jsonapi:"attr,surname"`
	UserID       uuid.UUID      `db:"user_id"`
	User         *User          `belongs_to:"user" jsonapi:"relation,user,omitempty"`
	Students     []*Student     `many_to_many:"parents_students" jsonapi:"relation,students,omitempty"`
	Appointments []*Appointment `many_to_many:"parents_appointments" jsonapi:"relation,appointments,omitempty"`
	Payments     []*Payment     `many_to_many:"parents_payments" jsonapi:"relation,payments,omitempty"`
}

// Return a string representation of the Parent resource
func (p Parent) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Parents is not required by pop and may be deleted
type Parents []*Parent

// Return a string representation of the Parents resource
func (p Parents) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*"
// (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *Parent) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{
			Field:   p.Email,
			Name:    "Email",
			Message: mailValidationMsg},
		&validators.StringIsPresent{
			Field:   p.Name,
			Name:    "Name",
			Message: nameValidationMsg},
		&validators.StringIsPresent{
			Field:   p.Surname,
			Name:    "Surname",
			Message: surnameValidationMsg},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Parent) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Parent) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a parent
func (p Parent) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/parents/%s", APIUrl, p.ID.String()),
	}
}

// JSONAPIRelationshipLinks is invoked for each relationship defined on the Parent struct when marshaled
func (p Parent) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "user" {
		return &jsonapi.Links{
			"user": fmt.Sprintf("http://%s/users/%s", APIUrl, p.UserID.String()),
		}
	}
	if relation == "students" {
		return &jsonapi.Links{
			"students": fmt.Sprintf("http://%s/parents/%s/students", APIUrl, p.ID.String()),
		}
	}
	if relation == "appointments" {
		return &jsonapi.Links{
			"appointments": fmt.Sprintf("http://%s/parents/%s/appointments", APIUrl, p.ID.String()),
		}
	}
	if relation == "payments" {
		return &jsonapi.Links{
			"payments": fmt.Sprintf("http://%s/parents/%s/payments", APIUrl, p.ID.String()),
		}
	}
	return nil
}
