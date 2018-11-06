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

	// Retrieve all Parents from the DB
	if err := q.All(parents); err != nil {
		return apiError(c, "Internal Error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Attatch users to parents
	for i, p := range *parents {
		user := &models.User{}
		if err := tx.Select("id", "created_at", "updated_at", "email", "role").
			Find(user, p.UserID); err != nil {
			return apiError(c, "Internal Error",
				"Internal Server Error", http.StatusInternalServerError, err)
		}
		p.User = user
		// Save it back to the parents list
		(*parents)[i] = p
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	// Convert the slice of parents to a slice of pointers to parents
	// because Pop wants the former, jsonapi, the latter
	parentsp := []*models.Parent{}
	for i := 0; i < len(*parents); i++ {
		parentsp = append(parentsp, &((*parents)[i]))
	}

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, parentsp)
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
	if err := tx.Eager("Students").Find(parent, c.Param("parent_id")); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// Convert the slice of students to a slice of pointers to students
	// because Pop wants the former, jsonapi, the latter
	studentsp := []*models.Student{}
	for i := 0; i < len(parent.Students); i++ {
		studentsp = append(studentsp, &(parent.Students[i]))
	}
	parent.StudentsRel = studentsp

	// Attatch the user to the parent
	user := &models.User{}
	if err := tx.Select("id", "created_at", "updated_at", "email", "role").
		Find(user, parent.UserID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	parent.User = user

	log.Println(parent)
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

	// Store the parent in the DB
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
