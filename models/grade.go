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

type Grade struct {
	ID        uuid.UUID `db:"id" jsonapi:"primary,grades"`
	CreatedAt time.Time `db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `db:"updated_at" jsonapi:"attr,created_at,iso8601"`
	Subject   string    `db:"subject" jsonapi:"attr,subject"`
	Grade     int       `db:"grade" jsonapi:"attr,grade"`
	StudentID uuid.UUID `db:"student_id"`
	Student   *Student  `belongs_to:"student" jsonapi:"relation,student,omitempty"`
}

// String is not required by pop and may be deleted
func (g Grade) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Grades is not required by pop and may be deleted
type Grades []*Grade

// String is not required by pop and may be deleted
func (g Grades) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (g *Grade) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (g *Grade) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (g *Grade) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a parent
func (grade Grade) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/grades/%s", APIUrl, grade.ID.String()),
	}
}

// Invoked for each relationship defined on the Grade struct when marshaled
func (grade Grade) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "student" {
		return &jsonapi.Links{
			"student": fmt.Sprintf("http://%s/students/%s", APIUrl, grade.StudentID.String()),
		}
	}
	return nil
}
