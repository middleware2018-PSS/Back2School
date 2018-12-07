package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createNotifications(tx *pop.Connection) error {
	johnDoeUser := &models.User{}
	paulaMillerUser := &models.User{}

	if err := tx.Find(johnDoeUser, johnDoe.UserID); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Find(paulaMillerUser, paulaMiller.UserID); err != nil {
		return errors.WithStack(err)
	}

	notification.Users = []*models.User{johnDoeUser, paulaMillerUser}

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
}
