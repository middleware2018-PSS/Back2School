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

			// Create parents
			if err := createParents(tx); err != nil {
				return err
			}

			// Create teachers
			if err := createTeachers(tx); err != nil {
				return err
			}

			// Create admins
			if err := createAdmins(tx); err != nil {
				return err
			}

			// Create students
			if err := createStudents(tx); err != nil {
				return err
			}

			// Create classes
			if err := createClasses(tx); err != nil {
				return err
			}

			// Create appointments
			if err := createAppointments(tx); err != nil {
				return err
			}

			return nil
		})
	})

})

func generateID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func parseDate(d string) time.Time {
	t, _ := time.Parse(time.RFC3339, d)
	return t
}
