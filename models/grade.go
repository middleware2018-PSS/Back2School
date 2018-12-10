package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cippaciong/jsonapi"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

// Grade is the model for a student's grade
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
func (g Grade) JSONAPILinks() *jsonapi.Links {
	return &jsonapi.Links{
		"self": fmt.Sprintf("http://%s/grades/%s", APIUrl, g.ID.String()),
	}
}

// JSONAPIRelationshipLinks is invoked for each relationship defined on the Grade struct when marshaled
func (g Grade) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "student" {
		return &jsonapi.Links{
			"student": fmt.Sprintf("http://%s/students/%s", APIUrl, g.StudentID.String()),
		}
	}
	return nil
}

// BelongsToParent implements the Ownable interface for grade/parent relationships
func (g Grade) BelongsToParent(tx *pop.Connection, pID string) bool {
	s := &Student{}
	if err := tx.Eager("Parents").Find(s, g.Student.ID); err != nil {
		log.Println("Error eager loading student parents")
		return false
	}
	log.Println(g.Student)
	for _, p := range s.Parents {
		if p.ID.String() == pID {
			return true
		}
	}
	return false
}

// BelongsToTeacher implements the Ownable interface for grade/teacher relationships
func (g Grade) BelongsToTeacher(tx *pop.Connection, tID string) bool {
	// Eager load the student's class
	s := &Student{}
	if err := tx.Eager("Class").Find(s, g.Student.ID); err != nil {
		log.Println("Error eager loading student class")
		return false
	}
	// Eager load the teachers of the student's class
	c := &Class{}
	if err := tx.Eager("Teachers").Find(c, s.Class.ID); err != nil {
		log.Println("Error eager loading class teachers")
		return false
	}
	for _, t := range c.Teachers {
		if t.ID.String() == tID {
			return true
		}
	}
	return false
}
