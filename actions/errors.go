package actions

import (
	"bytes"

	"github.com/cippaciong/jsonapi"
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

func apiError(c buffalo.Context, title, status string, httpcode int, err error) error {

	log.Debug("%+v", errors.WithStack(err))

	if ENV == "production" {
		res := new(bytes.Buffer)
		jsonapi.MarshalErrors(res, []*jsonapi.ErrorObject{{
			Title:  title,
			Detail: err.Error(),
			Status: status,
		}})
		return c.Render(httpcode,
			r.Func("application/json", customJSONRenderer(res.String())))
	} else {
		return errors.WithStack(err)
	}
}
