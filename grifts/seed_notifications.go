package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createNotifications(tx *pop.Connection) error {
	john_doe_user := &models.User{}
	paula_miller_user := &models.User{}

	if err := tx.Find(john_doe_user, john_doe.UserID); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Find(paula_miller_user, paula_miller.UserID); err != nil {
		return errors.WithStack(err)
	}

	notification.Users = []*models.User{john_doe_user, paula_miller_user}

	if _, err := tx.ValidateAndCreate(notification); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var notification *models.Notification = &models.Notification{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Message:   "Notification test",
	Time:      parseDate("2018-11-30T23:30:00Z"),
	//Users:     []*models.User{john_doe.User, paula_miller.User},
}
