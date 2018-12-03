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
	u_id := createUserFromTeacher(tx, *paula_miller)
	paula_miller.UserID = u_id
	if _, err := tx.ValidateAndCreate(paula_miller); err != nil {
		return errors.WithStack(err)
	}
	u_id = createUserFromTeacher(tx, *robert_smith)
	robert_smith.UserID = u_id
	if _, err := tx.ValidateAndCreate(robert_smith); err != nil {
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

var paula_miller *models.Teacher = &models.Teacher{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "paula.miller@example.com",
	Password:  "password",
	Name:      "Paula",
	Surname:   "Miller",
	Subject:   "Math",
}

var robert_smith *models.Teacher = &models.Teacher{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "robert.smith@example.com",
	Password:  "password",
	Name:      "Robert",
	Surname:   "Smith",
	Subject:   "English",
}
