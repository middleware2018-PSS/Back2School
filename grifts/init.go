package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/middleware2018-PSS/back2_school/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
