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
)

type Timetable struct {
	ID        uuid.UUID  `db:"id" jsonapi:"primary,timetables"`
	CreatedAt time.Time  `db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time  `db:"updated_at" jsonapi:"attr,created_at,iso8601"`
	Weekday   string     `db:"weekday" jsonapi:"attr,weekday"`
	Hour      string     `db:"hour" jsonapi:"attr,hour"`
	Subject   string     `db:"subject" jsonapi:"attr,subject"`
	ClassID   nulls.UUID `db:"class_id"`
	Class     *Class     `belongs_to:"class" jsonapi:"relation,class,omitempty"`
}

// String is not required by pop and may be deleted
func (t Timetable) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Timetables is not required by pop and may be deleted
type Timetables []*Timetable

// String is not required by pop and may be deleted
func (t Timetables) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Timetable) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Timetable) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Timetable) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// JSONAPILinks implements the Linkable interface for a timetable
func (t Timetable) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/timetables/%s", APIUrl, t.ID.String()),
	}
}

// JSONAPIRelationshipLinks is invoked for each relationship defined on the Timetable struct when marshaled
func (t Timetable) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "class" {
		return &jsonapi.Links{
			"class": fmt.Sprintf("http://%s/classes/%s",
				APIUrl, t.ClassID.UUID.String()),
		}
	}
	return nil
}

// BelongsToParent implements the Ownable interface for timetable/parent relationships
func (t Timetable) BelongsToParent(tx *pop.Connection, pID string) bool {
	return false
}

// BelongsToTeacher implements the Ownable interface for timetable/teacher relationships
func (t Timetable) BelongsToTeacher(tx *pop.Connection, tID string) bool {
	return true
	//for _, t := range c.Teachers {
	//if t.ID.String() == tID {
	//return true
	//}
	//}
	//return false
}
