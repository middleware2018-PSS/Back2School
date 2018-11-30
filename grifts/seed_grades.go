package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createGrades(tx *pop.Connection) error {
	if _, err := tx.ValidateAndCreate(grade); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var grade *models.Grade = &models.Grade{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Subject:   "math",
	Grade:     9,
	Student:   lisa_doe,
}
