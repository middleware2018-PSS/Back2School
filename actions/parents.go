package actions

import (
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
// Model: Singular (Parent)
// DB Table: Plural (parents)
// Resource: Plural (Parents)
// Path: Plural (/parents)
// View Template Folder: Plural (/templates/parents/)

// ParentsResource is the resource for the Parent model
type ParentsResource struct {
	buffalo.Resource
}

// List gets all Parents. This function is mapped to the path
// GET /parents
func (v ParentsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	parents := &models.Parents{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Parents from the DB
	if err := q.All(parents); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, parents))
}

// Show gets the data for one Parent. This function is mapped to
// the path GET /parents/{parent_id}
func (v ParentsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Parent
	parent := &models.Parent{}

	// To find the Parent the parameter parent_id is used.
	if err := tx.Find(parent, c.Param("parent_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, parent))
}

// New renders the form for creating a new Parent.
// This function is mapped to the path GET /parents/new
func (v ParentsResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Parent{}))
}

// Create adds a Parent to the DB. This function is mapped to the
// path POST /parents
func (v ParentsResource) Create(c buffalo.Context) error {
	// Allocate an empty Parent
	parent := &models.Parent{}

	// Bind parent to the html form elements
	if err := c.Bind(parent); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(parent)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, parent))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Parent was created successfully")

	// and redirect to the parents index page
	return c.Render(201, r.Auto(c, parent))
}

// Edit renders a edit form for a Parent. This function is
// mapped to the path GET /parents/{parent_id}/edit
func (v ParentsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Parent
	parent := &models.Parent{}

	if err := tx.Find(parent, c.Param("parent_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, parent))
}

// Update changes a Parent in the DB. This function is mapped to
// the path PUT /parents/{parent_id}
func (v ParentsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Parent
	parent := &models.Parent{}

	if err := tx.Find(parent, c.Param("parent_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Parent to the html form elements
	if err := c.Bind(parent); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(parent)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, parent))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Parent was updated successfully")

	// and redirect to the parents index page
	return c.Render(200, r.Auto(c, parent))
}

// Destroy deletes a Parent from the DB. This function is mapped
// to the path DELETE /parents/{parent_id}
func (v ParentsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Parent
	parent := &models.Parent{}

	// To find the Parent the parameter parent_id is used.
	if err := tx.Find(parent, c.Param("parent_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(parent); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Parent was destroyed successfully")

	// Redirect to the parents index page
	return c.Render(200, r.Auto(c, parent))
}