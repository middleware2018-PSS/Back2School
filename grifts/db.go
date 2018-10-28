package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/markbates/grift/grift"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			if err := tx.TruncateAll(); err != nil {
				return errors.WithStack(err)
			}

			// Create UUID
			id, _ := uuid.NewV4()

			user := &models.User{
				ID:        id,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Email:     "admin@example.com",
				Password:  "admin",
				Role:      "admin",
				Parent:    &models.Parent{},
				Teacher:   &models.Teacher{},
			}

			// Validate the data from the html form
			_, err := tx.ValidateAndCreate(user)
			if err != nil {
				return errors.WithStack(err)
			}

			return nil
		})
	})

})
