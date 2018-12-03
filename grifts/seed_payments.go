package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createPayments(tx *pop.Connection) error {
	if _, err := tx.ValidateAndCreate(payment); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var payment *models.Payment = &models.Payment{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	IssueDate: parseDate("2018-11-27T11:30:00Z"),
	DueDate:   parseDate("2019-11-27T11:30:00Z"),
	Amount:    1786.28,
	Parents:   []*models.Parent{john_doe, abbie_williams},
	Student:   lisa_doe,
}
