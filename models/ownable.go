package models

import "github.com/gobuffalo/pop"

type Ownable interface {
	BelongsToParent(tx *pop.Connection, pID string) bool
	BelongsToTeacher(tx *pop.Connection, tID string) bool
}
