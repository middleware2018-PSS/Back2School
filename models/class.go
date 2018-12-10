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

// Class is the model for a class of students
type Class struct {
	ID        uuid.UUID  `json:"id" db:"id" jsonapi:"primary,classes"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" jsonapi:"attr,created_at,iso8601"`
	Name      string     `json:"name" db:"name" jsonapi:"attr,name"`
	Year      time.Time  `json:"year" db:"year" jsonapi:"attr,year,iso8601"`
	Teachers  []*Teacher `many_to_many:"teachers_classes" jsonapi:"relation,teachers,omitempty"`
	Students  []*Student `has_many:"students" jsonapi:"relation,students,omitempty"`
}

// String is not required by pop and may be deleted
func (c Class) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Classes is not required by pop and may be deleted
type Classes []*Class

// String is not required by pop and may be deleted
func (c Classes) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Class) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Class) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Class) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a class
func (c Class) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/classes/%s", APIUrl, c.ID.String()),
	}
}

// JSONAPIRelationshipLinks is invoked for each relationship defined on the Class struct when marshaled
func (c Class) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "teachers" {
		return &jsonapi.Links{
			"teachers": fmt.Sprintf("http://%s/classes/%s/teachers",
				APIUrl, c.ID.String()),
		}
	}
	if relation == "students" {
		return &jsonapi.Links{
			"students": fmt.Sprintf("http://%s/classes/%s/students",
				APIUrl, c.ID.String()),
		}
	}
	return nil
}

// BelongsToParent implements the Ownable interface for class/parent relationships
func (c Class) BelongsToParent(tx *pop.Connection, pID string) bool {
	return false
}

// BelongsToTeacher implements the Ownable interface for class/teacher relationships
func (c Class) BelongsToTeacher(tx *pop.Connection, tID string) bool {
	for _, t := range c.Teachers {
		if t.ID.String() == tID {
			return true
		}
	}
	return false
}
