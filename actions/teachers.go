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
// Model: Singular (Teacher)
// DB Table: Plural (teachers)
// Resource: Plural (Teachers)
// Path: Plural (/teachers)
// View Template Folder: Plural (/templates/teachers/)

// TeachersResource is the resource for the Teacher model
type TeachersResource struct {
	buffalo.Resource
}

// List gets all Teachers. This function is mapped to the path
// GET /teachers
func (v TeachersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	teachers := &models.Teachers{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Teachers from the DB
	if err := q.All(teachers); err != nil {
		return apiError(c, "Internal Error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Attatch users to teachers
	for i, t := range *teachers {
		user := &models.User{}
		if err := tx.Select("id", "created_at", "updated_at", "email", "role").
			Find(user, t.UserID); err != nil {
			return apiError(c, "Internal Error",
				"Internal Server Error", http.StatusInternalServerError, err)
		}
		t.User = user
		// Save it back to the teachers list
		(*teachers)[i] = t
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	// Convert the slice of teachers to a slice of pointers to teachers
	// because Pop wants the former, jsonapi, the latter
	teachersp := []*models.Teacher{}
	for i := 0; i < len(*teachers); i++ {
		teachersp = append(teachersp, &((*teachers)[i]))
	}

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, teachersp)
	if err != nil {
		log.Debug("Problem marshalling teachers in actions.TeachersResource.List")
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Show gets the data for one Teacher. This function is mapped to
// the path GET /teachers/{teacher_id}
func (v TeachersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "No transaction found", "Internal Server Error",
			http.StatusInternalServerError, errors.New("No transaction found"))
	}

	// Allocate an empty Teacher
	teacher := &models.Teacher{}

	// To find the Teacher the parameter teacher_id is used.
	if err := tx.Find(teacher, c.Param("teacher_id")); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	// Attatch the user to the teacher
	user := &models.User{}
	if err := tx.Select("id", "created_at", "updated_at", "email", "role").
		Find(user, teacher.UserID); err != nil {
		return apiError(c, "The requested resource cannot be found",
			"Not Found", http.StatusNotFound, err)
	}

	teacher.User = user

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, teacher)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Create adds a Teacher to the DB. This function is mapped to the
// path POST /teachers
func (v TeachersResource) Create(c buffalo.Context) error {
	teacher := &models.Teacher{}

	// Unmarshall the JSON payload into a Teacher struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, teacher); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Create the User associated to the Teacher
	user := &models.User{
		Email:    teacher.Email,
		Password: teacher.Password,
		Role:     "teacher",
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

	log.Debug("User created in actions.TeachersResource.Create:\n%v\n", user)

	// Add the User ID to the Teacher
	teacher.UserID = user.ID

	// Store the teacher in the DB
	verrs, err = tx.ValidateAndCreate(teacher)
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
	teacher.Password = ""
	log.Debug("Teacher created in actions.TeachersResource.Create:\n%v\n", teacher)

	// If there are no errors return the Teacher resource
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, teacher)
	if err != nil {
		return apiError(c, "Error processing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Update changes a Teacher in the DB. This function is mapped to
// the path PUT /teachers/{teacher_id}
func (v TeachersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty Teacher
	teacher := &models.Teacher{}

	if err := tx.Find(teacher, c.Param("teacher_id")); err != nil {
		return apiError(c, "Cannot update the resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// Unmarshall the JSON payload into a Teacher struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, teacher); err != nil {
		return apiError(c, "Error processing the request payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	// Store the teacher in the DB
	verrs, err := tx.ValidateAndUpdate(teacher)
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
	teacher.Password = ""

	// Marshal the modified resource and send it back
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, teacher)
	if err != nil {
		return apiError(c, "Internal Error preparing the response payload",
			"Internal Server Error", http.StatusInternalServerError, err)
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// Destroy deletes a Teacher from the DB. This function is mapped
// to the path DELETE /teachers/{teacher_id}
func (v TeachersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty Teacher
	teacher := &models.Teacher{}

	// To find the Teacher the parameter teacher_id is used.
	if err := tx.Find(teacher, c.Param("teacher_id")); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// Allocate an empty User
	user := &models.User{}

	// Find the User with teacher.UserID
	if err := tx.Find(user, teacher.UserID); err != nil {
		return apiError(c, "Cannot delete resource. Resource not found",
			"Not Found", http.StatusNotFound, err)
	}

	// We delete only the user since the teacher entry is handled by cascading rules
	if err := tx.Destroy(user); err != nil {
		return apiError(c, "Internal error", "Internal Server Error",
			http.StatusInternalServerError, err)
	}

	// Redirect to the teachers index page
	return c.Render(204, r.Func("application/json",
		customJSONRenderer("")))
}
