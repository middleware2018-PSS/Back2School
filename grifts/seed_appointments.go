package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createAppointments(tx *pop.Connection) error {
	if _, err := tx.ValidateAndCreate(appointment); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var appointment *models.Appointment = &models.Appointment{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Time:      parseDate("2018-12-27T11:30:00Z"),
	Parents:   []*models.Parent{john_doe, abbie_williams},
	Teacher:   paula_miller,
	Student:   lisa_doe,
}
