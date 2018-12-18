package actions

import (
	"bytes"
	//"log"
	"net/http"

	"github.com/cippaciong/jsonapi"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/middleware2018-PSS/back2_school/models"
	"github.com/pkg/errors"
)

// ParentsResource is the resource for the Parent model
type ParentsResource struct {
	buffalo.Resource
}

// List gets all Parents. This function is mapped to the path
// GET /parents
// @Summary List parents
// @Description Get the list of all parents
// @Tags Parents
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Parent
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /parents [get]
// @Security ApiKeyAuth
func (v ParentsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	parents := &models.Parents{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	if c.Param("name") != "" {
		// Filter parents by name
		if err := q.Where("name = (?)", c.Param("name")).All(parents); err != nil {
			return apiError(c, "Internal Error", "Internal Server Error",
				http.StatusInternalServerError, err)
		}
	} else {
		// Retrieve all Parents from the DB
		if err := q.All(parents); err != nil {
			return apiError(c, "Internal Error", "Internal Server Error",
				http.StatusInternalServerError, err)
		}
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, *parents)
	if err != nil {
		log.Debug("Problem marshalling parents in actions.ParentsResource.List")
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Show gets the data for one Parent. This function is mapped to
// the path GET /parents/{parent_id}
// @Summary Get a parent
// @Description Get a single parent and its relationships
// @Tags Parents
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Parent
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /parents/{id} [get]
// @Security ApiKeyAuth
func (v ParentsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	// Allocate an empty Parent
	parent := &models.Parent{}

	// To find the Parent the parameter parent_id is used.
	if err := tx.Eager().Find(parent, c.Param("parent_id")); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// Attatch the user to the parent
	user := &models.User{}
	if err := tx.Select("id", "created_at", "updated_at", "email", "role").
		Find(user, parent.UserID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	parent.User = user

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, parent)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Create adds a Parent to the DB. This function is mapped to the
// path POST /parents
// @Summary Create a parent
// @Description Create a parent from the payload
// @Tags Parents
// @Accept  json
// @Produce  json
// @Param Parent body models.Parent true "Parent payload"
// @Success 200 {object} models.Parent
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 422 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /parents [post]
// @Security ApiKeyAuth
func (v ParentsResource) Create(c buffalo.Context) error {
	parent := &models.Parent{}

	// Unmarshall the JSON payload into a Parent struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, parent); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Create the User associated to the Parent
	user := &models.User{
		Email:    parent.Email,
		Password: parent.Password,
		Role:     "parent",
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Store the user in the DB
	verrs, err := tx.ValidateAndCreate(user)
	if err != nil {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, err)
	}

	// Check for validation errors
	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	log.Debug("User created in actions.ParentsResource.Create:\n%v\n", user)

	// Add the User ID to the Parent
	parent.UserID = user.ID

	// Store the parent in the DB
	verrs, err = tx.ValidateAndCreate(parent)
	if err != nil {
		return apiError(c, "Internal error",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Check for validation errors
	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Clear the Password so that it's not returned in the response
	parent.Password = ""
	log.Debug("Parent created in actions.ParentsResource.Create:\n%v\n", parent)

	// If there are no errors return the Parent resource
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, parent)
	if err != nil {
		return apiError(c, "Error processing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Update changes a Parent in the DB. This function is mapped to
// the path PUT /parents/{parent_id}
// @Summary Update a parent
// @Description Update a parent from the payload
// @Tags Parents
// @Accept  json
// @Produce  json
// @Param Parent body models.Parent true "Parent payload"
// @Success 200 {object} models.Parent
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 422 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /parents [put]
// @Security ApiKeyAuth
func (v ParentsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty Parent
	parent := &models.Parent{}

	if err := tx.Find(parent, c.Param("parent_id")); err != nil {
		return apiError(c, "Cannot update the resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// Unmarshall the JSON payload into a Parent struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, parent); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Update the parent in the DB
	verrs, err := tx.ValidateAndUpdate(parent)
	if err != nil {
		return apiError(c, "Internal error",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Check for validation errors
	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Nullify password before sending the response
	parent.Password = ""

	// Marshal the modified resource and send it back
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, parent)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Destroy deletes a Parent from the DB. This function is mapped
// to the path DELETE /parents/{parent_id}
// @Summary Delete a parent
// @Description Delete a parent
// @Tags Parents
// @Accept  json
// @Produce  json
// @Param  id path int true "Parent ID" Format(uuid)
// @Success 204 {object} models.Parent
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /parents/{id} [delete]
// @Security ApiKeyAuth
func (v ParentsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty Parent
	parent := &models.Parent{}

	// To find the Parent the parameter parent_id is used.
	if err := tx.Find(parent, c.Param("parent_id")); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// Allocate an empty User
	user := &models.User{}

	// Find the User with parent.UserID
	if err := tx.Find(user, parent.UserID); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// We delete only the user since the parent entry is handled by cascading rules
	if err := tx.Destroy(user); err != nil {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Redirect to the parents index page
	return c.Render(204, r.Func("application/json",
		customJSONRenderer("")))
}
