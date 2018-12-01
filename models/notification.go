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

type Notification struct {
	ID        uuid.UUID `db:"id" jsonapi:"primary,notifications"`
	CreatedAt time.Time `db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	// Attributes
	Time    time.Time `db:"time" jsonapi:"attr,time,iso8601"`
	Message string    `db:"message" jsonapi:"attr,message"`
	// Relationships
	Users []*User `many_to_many:"users_notifications" jsonapi:"relation,users,omitempty"`
}

// String is not required by pop and may be deleted
func (n Notification) String() string {
	jn, _ := json.Marshal(n)
	return string(jn)
}

// Notifications is not required by pop and may be deleted
type Notifications []*Notification

// String is not required by pop and may be deleted
func (n Notifications) String() string {
	jn, _ := json.Marshal(n)
	return string(jn)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (n *Notification) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (n *Notification) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (n *Notification) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a parent
func (notification Notification) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/notifications/%s",
			APIUrl, notification.ID.String()),
	}
}

// Invoked for each relationship defined on the Notification struct when marshaled
func (notification Notification) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "users" {
		return &jsonapi.Links{
			"users": fmt.Sprintf("http://%s/notifications/%s/users",
				APIUrl, notification.ID.String()),
		}
	}
	return nil
}
