package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cippaciong/jsonapi"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Student is the model for school students
type Student struct {
	ID           uuid.UUID      `json:"id" db:"id" jsonapi:"primary,students"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at" jsonapi:"attr,created_at,iso8601"`
	Name         string         `json:"name" db:"name" jsonapi:"attr,name"`
	Surname      string         `json:"surname" db:"surname" jsonapi:"attr,surname"`
	DateOfBirth  time.Time      `json:"date_of_birth" db:"date_of_birth" jsonapi:"attr,date_of_birth,iso8601"`
	Parents      []*Parent      `many_to_many:"parents_students" jsonapi:"relation,parents,omitempty"`
	Appointments []*Appointment `has_many:"appointments" jsonapi:"relation,appointments,omitempty"`
	ClassID      nulls.UUID     `db:"class_id"`
	Class        *Class         `belongs_to:"class" jsonapi:"relation,class,omitempty"`
	Grades       []*Grade       `has_many:"grades" jsonapi:"relation,grades,omitempty"`
	Payment      *Payment       `has_one:"payment" fk_id:"student_id" jsonapi:"relation,payment,omitempty"`
}

// String is not required by pop and may be deleted
func (s Student) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Students is not required by pop and may be deleted
type Students []*Student

// String is not required by pop and may be deleted
func (s Students) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Student) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{
			Field:   s.Name,
			Name:    "Name",
			Message: nameValidationMsg},
		&validators.StringIsPresent{
			Field:   s.Surname,
			Name:    "Surname",
			Message: surnameValidationMsg},
		&validators.TimeIsPresent{
			Field:   s.DateOfBirth,
			Name:    "DateOfBirth",
			Message: dateOfBirthValidationMsg},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Student) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Student) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a parent
func (s Student) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/students/%s", APIUrl, s.ID.String()),
	}
}

// JSONAPIRelationshipLinks is invoked for each relationship defined on the Student struct when marshaled
func (s Student) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "parents" {
		return &jsonapi.Links{
			"parents": fmt.Sprintf("http://%s/students/%s/parents", APIUrl, s.ID.String()),
		}
	}
	if relation == "grades" {
		return &jsonapi.Links{
			"grades": fmt.Sprintf("http://%s/students/%s/grades", APIUrl, s.ID.String()),
		}
	}
	if relation == "class" {
		return &jsonapi.Links{
			"class": fmt.Sprintf("http://%s/classes/%s", APIUrl, s.ClassID.UUID.String()),
		}
	}
	return nil
}

// BelongsToParent implements the Ownable interface for student/parent relationships
func (s Student) BelongsToParent(tx *pop.Connection, pID string) bool {
	for _, p := range s.Parents {
		if p.ID.String() == pID {
			return true
		}
	}
	return false
}

// BelongsToTeacher implements the Ownable interface for student/teacher relationships
func (s Student) BelongsToTeacher(tx *pop.Connection, tID string) bool {
	if !s.ClassID.Valid {
		return false
	}
	c := &Class{}
	if err := tx.Eager().Find(c, s.ClassID.UUID.String()); err != nil {
		return false
	}
	for _, t := range c.Teachers {
		if t.ID.String() == tID {
			return true
		}
	}
	return false
}
