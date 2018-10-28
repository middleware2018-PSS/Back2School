package actions

import (
	"io"

	"github.com/gobuffalo/buffalo/render"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		DefaultContentType: "application/json",
	})
}

func customJSONRenderer(payload string) func(io.Writer, render.Data) error {
	return func(w io.Writer, d render.Data) error {
		_, err := w.Write([]byte(payload))
		return err
	}
}
