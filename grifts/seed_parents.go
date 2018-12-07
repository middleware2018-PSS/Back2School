package grifts

import (
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createParents(tx *pop.Connection) error {
	uID := createUserFromParent(tx, *johnDoe)
	johnDoe.UserID = uID
	if _, err := tx.ValidateAndCreate(johnDoe); err != nil {
		return errors.WithStack(err)
	}
	uID = createUserFromParent(tx, *abbieWilliams)
	abbieWilliams.UserID = uID
	if _, err := tx.ValidateAndCreate(abbieWilliams); err != nil {
		return errors.WithStack(err)
	}

	return nil

}
func createUserFromParent(tx *pop.Connection, parent models.Parent) uuid.UUID {
	user := &models.User{
		Email:    parent.Email,
		Password: parent.Password,
		Role:     "parent",
	}

	// Store the user in the DB
	if _, err := tx.ValidateAndCreate(user); err != nil {
		log.Println(err)
	}
	return user.ID
}

var johnDoe *models.Parent = &models.Parent{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "john.doe@example.com",
	Password:  "password",
	Name:      "John",
	Surname:   "Doe",
}

var abbieWilliams *models.Parent = &models.Parent{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "abbie.williams@example.com",
	Password:  "password",
	Name:      "Abbie",
	Surname:   "Williams",
}
