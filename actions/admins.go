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
// Model: Singular (Admin)
// DB Table: Plural (admins)
// Resource: Plural (Admins)
// Path: Plural (/admins)
// View Template Folder: Plural (/templates/admins/)

// AdminsResource is the resource for the Admin model
type AdminsResource struct {
	buffalo.Resource
}

// List gets all Admins. This function is mapped to the path
// GET /admins
// @Summary List admins
// @Description Get the list of all admins
// @Tags Admins
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Admin
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /admins [get]
func (v AdminsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	admins := &models.Admins{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Admins from the DB
	if err := q.All(admins); err != nil {
		return apiError(c, "Internal Error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Attatch users to admins
	for i, a := range *admins {
		user := &models.User{}
		if err := tx.Select("id", "created_at", "updated_at", "email", "role").
			Find(user, a.UserID); err != nil {
			return apiError(c, "Internal Error",
				"Internal Server Error", http.StatusInternalServerError, err)
		}
		a.User = user
		// Save it back to the admins list
		(*admins)[i] = a
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	// Convert the slice of admins to a slice of pointers to admins
	// because Pop wants the former, jsonapi, the latter
	adminsp := []*models.Admin{}
	for i := 0; i < len(*admins); i++ {
		adminsp = append(adminsp, &((*admins)[i]))
	}

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, adminsp)
	if err != nil {
		log.Debug("Problem marshalling admins in actions.AdminsResource.List")
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Show gets the data for one Admin. This function is mapped to
// the path GET /admins/{admin_id}
// @Summary Get an admin
// @Description Get a single admin and its relationships
// @Tags Admins
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Admin
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /admins/{id} [get]
func (v AdminsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	// Allocate an empty Admin
	admin := &models.Admin{}

	// To find the Admin the parameter admin_id is used.
	if err := tx.Find(admin, c.Param("admin_id")); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// Attatch the user to the admin
	user := &models.User{}
	if err := tx.Select("id", "created_at", "updated_at", "email", "role").
		Find(user, admin.UserID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	admin.User = user

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, admin)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Create adds a Admin to the DB. This function is mapped to the
// path POST /admins
// @Summary Create an admin
// @Description Create an admin from the payload
// @Tags Admins
// @Accept  json
// @Produce  json
// @Param Admin body models.Admin true "Admin payload"
// @Success 200 {object} models.Admin
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 422 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /admins [post]
func (v AdminsResource) Create(c buffalo.Context) error {
	admin := &models.Admin{}

	// Unmarshall the JSON payload into a Admin struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, admin); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Create the User associated to the Admin
	user := &models.User{
		Email:    admin.Email,
		Password: admin.Password,
		Role:     "admin",
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

	log.Debug("User created in actions.AdminsResource.Create:\n%v\n", user)

	// Add the User ID to the Admin
	admin.UserID = user.ID

	// Store the admin in the DB
	verrs, err = tx.ValidateAndCreate(admin)
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
	admin.Password = ""
	log.Debug("Admin created in actions.AdminsResource.Create:\n%v\n", admin)

	// If there are no errors return the Admin resource
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, admin)
	if err != nil {
		return apiError(c, "Error processing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Update changes a Admin in the DB. This function is mapped to
// the path PUT /admins/{admin_id}
// @Summary Update an admin
// @Description Update an admin from the payload
// @Tags Admins
// @Accept  json
// @Produce  json
// @Param Admin body models.Admin true "Admin payload"
// @Success 200 {object} models.Admin
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 422 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /admins [put]
func (v AdminsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty Admin
	admin := &models.Admin{}

	if err := tx.Find(admin, c.Param("admin_id")); err != nil {
		return apiError(c, "Cannot update the resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// Unmarshall the JSON payload into a Admin struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, admin); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Store the admin in the DB
	verrs, err := tx.ValidateAndUpdate(admin)
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
	admin.Password = ""

	// Marshal the modified resource and send it back
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, admin)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Destroy deletes a Admin from the DB. This function is mapped
// to the path DELETE /admins/{admin_id}
// @Summary Delete an admin
// @Description Delete an admin
// @Tags Admins
// @Accept  json
// @Produce  json
// @Param  id path int true "Admin ID" Format(uuid)
// @Success 204 {object} models.Admin
// @Failure 404 {object} jsonapi.ErrorObject
// @Failure 500 {object} jsonapi.ErrorObject
// @Router /admins/{id} [delete]
func (v AdminsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty Admin
	admin := &models.Admin{}

	// To find the Admin the parameter admin_id is used.
	if err := tx.Find(admin, c.Param("admin_id")); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// Allocate an empty User
	user := &models.User{}

	// Find the User with admin.UserID
	if err := tx.Find(user, admin.UserID); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// We delete only the user since the admin entry is handled by cascading rules
	if err := tx.Destroy(user); err != nil {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Redirect to the admins index page
	return c.Render(204, r.Func("application/json",
		customJSONRenderer("")))
}
