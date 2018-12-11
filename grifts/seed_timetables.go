package grifts

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

func createTimetables(tx *pop.Connection) error {

	hours := []string{"09", "10", "11", "12", "14", "15", "16"}
	subjects := []string{"Math", "Science", "History", "Geography", "Art", "Music", "English"}
	for i, h := range hours {
		timetable.ID = generateID()
		timetable.Hour = h
		timetable.Subject = subjects[i]
		if _, err := tx.ValidateAndCreate(timetable); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

var timetable *models.Timetable = &models.Timetable{
	ID:        generateID(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Weekday:   "Monday",
	Hour:      "09",
	Subject:   "Math",
	ClassID:   nulls.NewUUID(class1a.ID),
}
