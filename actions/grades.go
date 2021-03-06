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
// Model: Singular (Grade)
// DB Table: Plural (grades)
// Resource: Plural (Grades)
// Path: Plural (/grades)
// View Template Folder: Plural (/templates/grades/)

// GradesResource is the resource for the Grade model
type GradesResource struct {
	buffalo.Resource
}

// List gets all Grades. This function is mapped to the path
// GET /grades
// @Summary List grades
// @Description Get the list of all grades
// @Tags Grades
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Grade
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /grades [get]
// @Security ApiKeyAuth
func (v GradesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	grades := &models.Grades{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Grades from the DB
	if err := q.All(grades); err != nil {
		return apiError(c, "Internal Error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, *grades)
	if err != nil {
		log.Debug("Problem marshalling grades in actions.GradesResource.List")
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Show gets the data for one Grade. This function is mapped to
// the path GET /grades/{grade_id}
// @Summary Get a grade
// @Description Get a single grade and its relationships
// @Tags Grades
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Grade
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /grades/{id} [get]
// @Security ApiKeyAuth
func (v GradesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	// Allocate an empty Grade
	grade := &models.Grade{}

	// To find the Grade the parameter grade_id is used.
	if err := tx.Eager().Find(grade, c.Param("grade_id")); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, grade)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// New renders the form for creating a new Grade.
// This function is mapped to the path GET /grades/new
func (v GradesResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Grade{}))
}

// Create adds a Grade to the DB. This function is mapped to the
// path POST /grades
// @Summary Create a grade
// @Description Create a grade from the payload
// @Tags Grades
// @Accept  json
// @Produce  json
// @Param Grade body models.Grade true "Grade payload"
// @Success 200 {object} models.Grade
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 422 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /grades [post]
// @Security ApiKeyAuth
func (v GradesResource) Create(c buffalo.Context) error {
	// Allocate an empty Grade
	grade := &models.Grade{}

	// Unmarshal grade from the json payload
	if err := jsonapi.UnmarshalPayload(c.Request().Body, grade); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Create and save the grade
	verrs, err := tx.ValidateAndCreate(grade)
	if err != nil {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, err)
	}

	// Check for validation errors
	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Log grade creation
	log.Debug("Grade created in actions.GradesResource.Create:\n%v\n", grade)

	// Reload the grade to rebuild relationships
	if err := tx.Eager().Find(grade, grade.ID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// If there are no errors return the Grade resource
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, grade)
	if err != nil {
		return apiError(c, "Error processing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}
	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Edit renders a edit form for a Grade. This function is
// mapped to the path GET /grades/{grade_id}/edit
func (v GradesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Grade
	grade := &models.Grade{}

	if err := tx.Find(grade, c.Param("grade_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, grade))
}

// Update changes a Grade in the DB. This function is mapped to
// the path PUT /grades/{grade_id}
// @Summary Update a grade
// @Description Update a grade from the payload
// @Tags Grades
// @Accept  json
// @Produce  json
// @Param Grade body models.Grade true "Grade payload"
// @Success 200 {object} models.Grade
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 422 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /grades [put]
// @Security ApiKeyAuth
func (v GradesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty Grade
	grade := &models.Grade{}

	if err := tx.Find(grade, c.Param("grade_id")); err != nil {
		return apiError(c, "Cannot update the resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// Unmarshall the JSON payload into a Grade struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, grade); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Update the grade in the DB
	verrs, err := tx.ValidateAndUpdate(grade)
	if err != nil {
		return apiError(c, "Internal error",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Check for validation errors
	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Reload the grade to rebuild relationships
	if err := tx.Eager().Find(grade, grade.ID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// Marshal the resource and send it back
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, grade)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Destroy deletes a Grade from the DB. This function is mapped
// to the path DELETE /grades/{grade_id}
// @Summary Delete a grade
// @Description Delete a grade
// @Tags Grades
// @Accept  json
// @Produce  json
// @Param  id path int true "Grade ID" Format(uuid)
// @Success 204 {object} models.Grade
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /grades/{id} [delete]
// @Security ApiKeyAuth
func (v GradesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty Grade
	grade := &models.Grade{}

	// To find the Grade the parameter grade_id is used.
	if err := tx.Find(grade, c.Param("grade_id")); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	if err := tx.Destroy(grade); err != nil {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Redirect to the grades index page
	return c.Render(204, r.Func("application/json",
		customJSONRenderer("")))
}
