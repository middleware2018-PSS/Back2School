package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createClasses(tx *pop.Connection) error {
	if _, err := tx.ValidateAndCreate(class_1a); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var class_1a *models.Class = &models.Class{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Name:      "1A",
	Year:      parseDate("2018-09-01T00:00:00Z"),
	Teachers:  []*models.Teacher{paula_miller, robert_smith},
	Students:  []*models.Student{lisa_doe},
}
