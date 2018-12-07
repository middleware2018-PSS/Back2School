package grifts

import (
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createAdmins(tx *pop.Connection) error {
	uID := createUserFromAdmin(tx, *angelaKennedy)
	angelaKennedy.UserID = uID
	if _, err := tx.ValidateAndCreate(angelaKennedy); err != nil {
		return errors.WithStack(err)
	}
	uID = createUserFromAdmin(tx, *chrisWhite)
	chrisWhite.UserID = uID
	if _, err := tx.ValidateAndCreate(chrisWhite); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func createUserFromAdmin(tx *pop.Connection, admin models.Admin) uuid.UUID {
	user := &models.User{
		Email:    admin.Email,
		Password: admin.Password,
		Role:     "admin",
	}

	// Store the user in the DB
	if _, err := tx.ValidateAndCreate(user); err != nil {
		log.Println(err)
	}
	return user.ID
}

var angelaKennedy *models.Admin = &models.Admin{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "angela.kennedy@example.com",
	Password:  "password",
	Name:      "Angela",
	Surname:   "Kennedy",
}

var chrisWhite *models.Admin = &models.Admin{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "chris.white@example.com",
	Password:  "password",
	Name:      "Chris",
	Surname:   "White",
}
