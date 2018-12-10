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
// Model: Singular (User)
// DB Table: Plural (users)
// Resource: Plural (Users)
// Path: Plural (/users)
// View Template Folder: Plural (/templates/users/)

// UsersResource is the resource for the User model
type UsersResource struct {
	buffalo.Resource
}

// List gets all Users. This function is mapped to the path
// GET /users
// @Summary List users
// @Description Get the list of all users
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /users [get]
// @Security ApiKeyAuth
func (v UsersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	users := &models.Users{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Users from the DB
	if err := q.Select("id", "created_at", "updated_at", "email", "role").
		All(users); err != nil {
		return apiError(c, "Internal Error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, *users)
	if err != nil {
		log.Debug("Problem marshalling users in actions.UsersResource.List")
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Show gets the data for one User. This function is mapped to
// the path GET /users/{user_id}
// @Summary Get a user
// @Description Get a single user and its relationships
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /users/{id} [get]
// @Security ApiKeyAuth
func (v UsersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	// We omit the Password column because we don't want to return that
	if err := tx.Eager("Notifications").Select("id", "created_at", "updated_at", "email", "role").
		Find(user, c.Param("user_id")); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, user)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	log.Println(user)

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Create adds a User to the DB. This function is mapped to the
// path POST /users
// @Summary Create a user
// @Description Create a user from the payload
// @Tags Users
// @Accept  json
// @Produce  json
// @Param User body models.User true "User payload"
// @Success 200 {object} models.User
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 422 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /users [post]
// @Security ApiKeyAuth
func (v UsersResource) Create(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}

	// Unmarshall the JSON payload into a User struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, user); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Validate the data from the payload
	verrs, err := tx.ValidateAndCreate(user)
	if err != nil {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, err)
	}

	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Log user creation
	log.Debug("User created in actions.UsersResource.Create:\n%v\n", user)

	// Reload the user to rebuild relationships
	if err := tx.Eager("Notifications").
		Select("id", "created_at", "updated_at", "email", "role").
		Find(user, user.ID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// If there are no errors return the User resource
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, user)
	if err != nil {
		return apiError(c, "Error processing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}
	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))

}

// Edit renders a edit form for a User. This function is
// mapped to the path GET /users/{user_id}/edit
func (v UsersResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, user))
}

// Update changes a User in the DB. This function is mapped to
// the path PUT /users/{user_id}
// @Summary Update a user
// @Description Update a user from the payload
// @Tags Users
// @Accept  json
// @Produce  json
// @Param User body models.User true "User payload"
// @Success 200 {object} models.User
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 422 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /users [put]
// @Security ApiKeyAuth
func (v UsersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return apiError(c, "Cannot update the resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// Unmarshall the JSON payload into a User struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, user); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Update the user in the DB
	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return apiError(c, "Internal error",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Check for validation errors
	if verrs.HasAny() {
		return apiError(c, "Validation Error", "Unprocessable Entity",
			http.StatusUnprocessableEntity, verrs)
	}

	// Reload the user to rebuild relationships
	if err := tx.Eager("Notifications").
		Select("id", "created_at", "updated_at", "email", "role").
		Find(user, user.ID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// Marshal the resource and send it back
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, user)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Destroy deletes a User from the DB. This function is mapped
// to the path DELETE /users/{user_id}
// @Summary Delete a user
// @Description Delete a user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param  id path int true "User ID" Format(uuid)
// @Success 204 {object} models.User
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /users/{id} [delete]
// @Security ApiKeyAuth
func (v UsersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	if err := tx.Destroy(user); err != nil {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Redirect to the parents index page
	return c.Render(204, r.Func("application/json",
		customJSONRenderer("")))
}
