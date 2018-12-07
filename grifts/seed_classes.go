package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createClasses(tx *pop.Connection) error {
	if _, err := tx.ValidateAndCreate(class1a); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var class1a *models.Class = &models.Class{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Name:      "1A",
	Year:      parseDate("2018-09-01T00:00:00Z"),
	Teachers:  []*models.Teacher{paulaMiller, robertSmith},
	Students:  []*models.Student{lisaDoe},
}
