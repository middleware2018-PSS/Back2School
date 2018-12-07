package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createStudents(tx *pop.Connection) error {
	if _, err := tx.ValidateAndCreate(lisaDoe); err != nil {
		return errors.WithStack(err)
	}
	if _, err := tx.ValidateAndCreate(alexDoe); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var lisaDoe *models.Student = &models.Student{
	ID:          generateID(),
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
	Name:        "Lisa",
	Surname:     "Doe",
	DateOfBirth: parseDate("2010-03-27T00:00:00Z"),
	Parents:     []*models.Parent{johnDoe, abbieWilliams},
}

var alexDoe *models.Student = &models.Student{
	ID:          generateID(),
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
	Name:        "Alex",
	Surname:     "Doe",
	DateOfBirth: parseDate("2012-07-02T00:00:00Z"),
	Parents:     []*models.Parent{johnDoe, abbieWilliams},
}
