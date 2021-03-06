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

// Teacher is the model for school teachers
type Teacher struct {
	ID           uuid.UUID      `json:"id" db:"id" jsonapi:"primary,teachers"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	Email        string         `json:"email" db:"email" jsonapi:"attr,email"`
	Password     string         `json:"passowrd" db:"-" jsonapi:"attr,password,omitempty"`
	Name         string         `json:"name" db:"name" jsonapi:"attr,name"`
	Surname      string         `json:"surname" db:"surname" jsonapi:"attr,surname"`
	Subject      string         `json:"subject" db:"subject" jsonapi:"attr,subject"`
	UserID       uuid.UUID      `db:"user_id"`
	User         *User          `belongs_to:"user" jsonapi:"relation,user,omitempty"`
	Appointments []*Appointment `has_many:"appointments" jsonapi:"relation,appointments,omitempty"`
	Classes      []*Class       `many_to_many:"teachers_classes" jsonapi:"relation,classes,omitempty"`
}

// String is not required by pop and may be deleted
func (t Teacher) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Teachers is not required by pop and may be deleted
type Teachers []*Teacher

// String is not required by pop and may be deleted
func (t Teachers) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*"
// (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *Teacher) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Field: t.Email, Name: "Email", Message: mailValidationMsg},
		&validators.StringIsPresent{Field: t.Name, Name: "Name", Message: nameValidationMsg},
		&validators.StringIsPresent{Field: t.Surname, Name: "Surname", Message: surnameValidationMsg},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Teacher) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Teacher) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a teacher
func (t Teacher) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/teachers/%s", APIUrl, t.ID.String()),
	}
}

// JSONAPIRelationshipLinks is invoked for each relationship defined on the Teacher struct when marshaled
func (t Teacher) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "user" {
		return &jsonapi.Links{
			"user": fmt.Sprintf("http://%s/users/%s", APIUrl, t.UserID.String()),
		}
	}
	if relation == "classes" {
		return &jsonapi.Links{
			"classes": fmt.Sprintf("http://%s/teachers/%s/classes", APIUrl, t.ID.String()),
		}
	}
	if relation == "appointments" {
		return &jsonapi.Links{
			"appointments": fmt.Sprintf("http://%s/teachers/%s/appointments", APIUrl, t.ID.String()),
		}
	}
	return nil
}
