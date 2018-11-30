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
	u_id := createUserFromAdmin(tx, *angela_kennedy)
	angela_kennedy.UserID = u_id
	if _, err := tx.ValidateAndCreate(angela_kennedy); err != nil {
		return errors.WithStack(err)
	}
	u_id = createUserFromAdmin(tx, *chris_white)
	chris_white.UserID = u_id
	if _, err := tx.ValidateAndCreate(chris_white); err != nil {
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

var angela_kennedy *models.Admin = &models.Admin{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "angela.kennedy@example.com",
	Password:  "password",
	Name:      "Angela",
	Surname:   "Kennedy",
}

var chris_white *models.Admin = &models.Admin{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "chris.white@example.com",
	Password:  "password",
	Name:      "Chris",
	Surname:   "White",
}
