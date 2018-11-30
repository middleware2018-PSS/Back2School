package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createStudents(tx *pop.Connection) error {
	if _, err := tx.ValidateAndCreate(lisa_doe); err != nil {
		return errors.WithStack(err)
	}
	if _, err := tx.ValidateAndCreate(alex_doe); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var lisa_doe *models.Student = &models.Student{
	ID:          generateID(),
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
	Name:        "Lisa",
	Surname:     "Doe",
	DateOfBirth: parseDate("2010-03-27T00:00:00Z"),
	Parents:     []*models.Parent{john_doe, abbie_williams},
}

var alex_doe *models.Student = &models.Student{
	ID:          generateID(),
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
	Name:        "Alex",
	Surname:     "Doe",
	DateOfBirth: parseDate("2012-07-02T00:00:00Z"),
	Parents:     []*models.Parent{john_doe, abbie_williams},
}
