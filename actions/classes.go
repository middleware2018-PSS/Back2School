package actions

import (
	"bytes"
	"net/http"

	"github.com/cippaciong/jsonapi"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Class)
// DB Table: Plural (classes)
// Resource: Plural (Classes)
// Path: Plural (/classes)
// View Template Folder: Plural (/templates/classes/)

// ClassesResource is the resource for the Class model
type ClassesResource struct {
	buffalo.Resource
}

// List gets all Classes. This function is mapped to the path
// GET /classes
func (v ClassesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	classes := &models.Classes{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Classes from the DB
	if err := q.All(classes); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, classes))
}

// Show gets the data for one Class. This function is mapped to
// the path GET /classes/{class_id}
func (v ClassesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Class
	class := &models.Class{}

	// To find the Class the parameter class_id is used.
	if err := tx.Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, class))
}

// New renders the form for creating a new Class.
// This function is mapped to the path GET /classes/new
func (v ClassesResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Class{}))
}

// Create adds a Class to the DB. This function is mapped to the
// path POST /classes
func (v ClassesResource) Create(c buffalo.Context) error {
	// Allocate an empty Class
	class := &models.Class{}

	// Unmarshal class from the json payload
	if err := jsonapi.UnmarshalPayload(c.Request().Body, class); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Create and save the class
	verrs, err := tx.ValidateAndCreate(class)
	if err != nil {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, err)
	}

	// Check for validation errors
	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Log class creation
	log.Debug("Class created in actions.ClasssesResource.Create:\n%v\n", class)

	// Reload the class to rebuild relationships
	if err := tx.Eager().Find(class, class.ID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// If there are no errors return the Appointment resource
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, class)
	if err != nil {
		return apiError(c, "Error processing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}
	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Edit renders a edit form for a Class. This function is
// mapped to the path GET /classes/{class_id}/edit
func (v ClassesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Class
	class := &models.Class{}

	if err := tx.Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, class))
}

// Update changes a Class in the DB. This function is mapped to
// the path PUT /classes/{class_id}
func (v ClassesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Class
	class := &models.Class{}

	if err := tx.Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Class to the html form elements
	if err := c.Bind(class); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(class)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, class))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Class was updated successfully")

	// and redirect to the classes index page
	return c.Render(200, r.Auto(c, class))
}

// Destroy deletes a Class from the DB. This function is mapped
// to the path DELETE /classes/{class_id}
func (v ClassesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Class
	class := &models.Class{}

	// To find the Class the parameter class_id is used.
	if err := tx.Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(class); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Class was destroyed successfully")

	// Redirect to the classes index page
	return c.Render(200, r.Auto(c, class))
}
