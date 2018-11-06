package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Student struct {
	ID          uuid.UUID `json:"id" db:"id" jsonapi:"primary,students"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,created_at,iso8601"`
	Name        string    `json:"name" db:"name" jsonapi:"attr,name"`
	Surname     string    `json:"surname" db:"surname" jsonapi:"attr,surname"`
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth" jsonapi:"attr,date_of_birth,iso8601"`
}

// String is not required by pop and may be deleted
func (s Student) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Students is not required by pop and may be deleted
type Students []Student

// String is not required by pop and may be deleted
func (s Students) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Student) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
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
