package grifts

import (
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createTeachers(tx *pop.Connection) error {
	uID := createUserFromTeacher(tx, *paulaMiller)
	paulaMiller.UserID = uID
	if _, err := tx.ValidateAndCreate(paulaMiller); err != nil {
		return errors.WithStack(err)
	}
	uID = createUserFromTeacher(tx, *robertSmith)
	robertSmith.UserID = uID
	if _, err := tx.ValidateAndCreate(robertSmith); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
func createUserFromTeacher(tx *pop.Connection, teacher models.Teacher) uuid.UUID {
	user := &models.User{
		Email:    teacher.Email,
		Password: teacher.Password,
		Role:     "teacher",
	}

	// Store the user in the DB
	if _, err := tx.ValidateAndCreate(user); err != nil {
		log.Println(err)
	}
	return user.ID
}

var paulaMiller *models.Teacher = &models.Teacher{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "paula.miller@example.com",
	Password:  "password",
	Name:      "Paula",
	Surname:   "Miller",
	Subject:   "Math",
}

var robertSmith *models.Teacher = &models.Teacher{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "robert.smith@example.com",
	Password:  "password",
	Name:      "Robert",
	Surname:   "Smith",
	Subject:   "English",
}
