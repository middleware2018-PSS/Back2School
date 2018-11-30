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
	u_id := createUserFromParent(tx, *john_doe)
	john_doe.UserID = u_id
	if _, err := tx.ValidateAndCreate(john_doe); err != nil {
		return errors.WithStack(err)
	}
	u_id = createUserFromParent(tx, *abbie_williams)
	abbie_williams.UserID = u_id
	if _, err := tx.ValidateAndCreate(abbie_williams); err != nil {
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

var john_doe *models.Parent = &models.Parent{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "john.doe@example.com",
	Password:  "password",
	Name:      "John",
	Surname:   "Doe",
}

var abbie_williams *models.Parent = &models.Parent{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Email:     "abbie.williams@example.com",
	Password:  "password",
	Name:      "Abbie",
	Surname:   "Williams",
}
