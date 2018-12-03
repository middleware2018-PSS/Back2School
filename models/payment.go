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

type Payment struct {
	ID        uuid.UUID `json:"id" db:"id" jsonapi:"primary,payments"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	// Attributes
	DueDate   time.Time `json:"due_date" db:"due_date" jsonapi:"attr,due_date,iso8601"`
	IssueDate time.Time `json:"issue_date" db:"issue_date" jsonapi:"attr,issue_date,iso8601"`
	Amount    float64   `json:"amount" db:"amount" jsonapi:"attr,amount"`
	// Relationships
	StudentID uuid.UUID `db:"student_id"`
	Student   *Student  `belongs_to:"student" jsonapi:"relation,student,omitempty"`
	Parents   []*Parent `many_to_many:"parents_payments" jsonapi:"relation,parents,omitempty"`
}

// String is not required by pop and may be deleted
func (p Payment) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Payments is not required by pop and may be deleted
type Payments []*Payment

// String is not required by pop and may be deleted
func (p Payments) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Payment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Payment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Payment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a payment
func (payment Payment) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/payments/%s", APIUrl, payment.ID.String()),
	}
}

// Invoked for each relationship defined on the Payment struct when marshaled
func (payment Payment) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "parents" {
		return &jsonapi.Links{
			"parents": fmt.Sprintf("http://%s/payments/%s/parents", APIUrl, payment.ID.String()),
		}
	}
	if relation == "student" {
		return &jsonapi.Links{
			"student": fmt.Sprintf("http://%s/students/%s", APIUrl, payment.StudentID.String()),
		}
	}
	return nil
}
