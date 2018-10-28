package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Parent struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Surname   string    `json:"surname" db:"surname"`
	UserID    uuid.UUID `db:"user_id"`
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
		&validators.EmailIsPresent{Field: p.Email, Name: "Email", Message: "Mail is not in the right format."},
		&validators.StringIsPresent{Field: p.Name, Name: "Name"},
		&validators.StringIsPresent{Field: p.Surname, Name: "Surname"},
	), nil
}
