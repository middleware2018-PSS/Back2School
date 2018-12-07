package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cippaciong/jsonapi"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Appointment struct {
	ID        uuid.UUID `json:"id" db:"id" jsonapi:"primary,appointments"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	Time      time.Time `json:"time" db:"time" jsonapi:"attr,time,iso8601"`
	Teacher   *Teacher  `belongs_to:"teacher" jsonapi:"relation,teacher,omitempty"`
	TeacherID uuid.UUID `db:"teacher_id"`
	Parents   []*Parent `many_to_many:"parents_appointments" jsonapi:"relation,parents,omitempty"`
	Student   *Student  `belongs_to:"student" jsonapi:"relation,student,omitempty"`
	StudentID uuid.UUID `db:"student_id"`
}

// String is not required by pop and may be deleted
func (a Appointment) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Appointments is not required by pop and may be deleted
type Appointments []*Appointment

// String is not required by pop and may be deleted
func (a Appointments) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Appointment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Appointment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Appointment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a parent
func (appointment Appointment) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/appointments/%s", APIUrl, appointment.ID.String()),
	}
}

// Invoked for each relationship defined on the Appointment struct when marshaled
func (appointment Appointment) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "student" {
		return &jsonapi.Links{
			"student": fmt.Sprintf("http://%s/students/%s", APIUrl, appointment.Student.ID.String()),
		}
	}
	if relation == "parents" {
		return &jsonapi.Links{
			"parents": fmt.Sprintf("http://%s/appointments/%s/parents", APIUrl, appointment.ID.String()),
		}
	}
	if relation == "teacher" {
		return &jsonapi.Links{
			"teacher": fmt.Sprintf("http://%s/teacher/%s", APIUrl, appointment.Teacher.ID.String()),
		}
	}
	return nil
}
