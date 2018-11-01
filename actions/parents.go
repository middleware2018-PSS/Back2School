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
		return apiError(c, "no transaction found", "InternalServerError",
			http.StatusInternalServerError, errors.New("no transaction found"))
	}

	// Allocate an empty slice of Parents
	parents := &models.Parents{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Parents from the DB
	if err := q.All(parents); err != nil {
		//TODO 404
		log.Println("HERE")
		return apiError(c, "Error retrieving parents from the DB", "InternalServerError",
			http.StatusInternalServerError, err)
		//return errors.WithStack(err)
	}

	// Attatch users to parents
	for i, p := range *parents {
		user := &models.User{}
		if err := tx.Select("id", "created_at", "updated_at", "email", "role").
			Find(user, p.UserID); err != nil {
			return apiError(c, "Cannot find parent resource with the specified id",
				"Not Found", http.StatusNotFound, err)
			//return c.Error(404, err)
		}
		p.User = user
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
		log.Println("Problem marshalling parents")
		return c.Render(http.StatusInternalServerError, r.JSON(err.Error()))
	}

	return c.Render(200, r.Func("application/json", customJSONRenderer(res.String())))
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
		return apiError(c, "Error retrieving parents from the DB", "InternalServerError",
			http.StatusInternalServerError, err)
		//return c.Error(404, err)
	}

	// Attatch the user to the parent
	user := &models.User{}
	if err := tx.Select("id", "created_at", "updated_at", "email", "role").
		Find(user, parent.UserID); err != nil {
		return c.Error(404, err)
	}

	parent.User = user

	res := new(bytes.Buffer)
	err := jsonapi.MarshalPayload(res, parent)
	if err != nil {
		log.Println("Problem marshalling the Parent in Show()")
		return c.Render(http.StatusInternalServerError, r.JSON(err.Error()))
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
}

// New renders the form for creating a new Parent.
// This function is mapped to the path GET /parents/new
func (v ParentsResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Parent{}))
}

// Create adds a Parent to the DB. This function is mapped to the
// path POST /parents
func (v ParentsResource) Create(c buffalo.Context) error {
	parent := &models.Parent{}

	// Unmarshall the JSON payload into a Parent struct
	if err := jsonapi.UnmarshalPayload(c.Request().Body, parent); err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Store the user in the DB
	verrs, err := tx.ValidateAndCreate(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		if ENV == "production" {
			res := new(bytes.Buffer)
			jsonapi.MarshalErrors(res, []*jsonapi.ErrorObject{{
				Title:  "Validation Error",
				Detail: verrs.Error(),
				Status: "Unprocessable Entity",
			}})
			return c.Render(http.StatusUnprocessableEntity,
				r.Func("application/json", customJSONRenderer(res.String())))
		} else {
			return errors.WithStack(verrs)
		}
	}

	log.Printf("PRINT USER CREATED IN PARENT: %v", user)

	// Add the User ID to the Parent
	parent.UserID = user.ID

	// Store the parent in the DB
	verrs, err = tx.ValidateAndCreate(parent)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		res := new(bytes.Buffer)
		jsonapi.MarshalErrors(res, []*jsonapi.ErrorObject{{
			Title:  "Validation Error",
			Detail: verrs.Error(),
			Status: "Unprocessable Entity",
		}})
		return c.Render(http.StatusUnprocessableEntity, r.Func("application/json",
			customJSONRenderer(res.String())))
	}

	// Clear the Password so that it's not returned in the response
	parent.Password = ""
	log.Printf("PRINT PARENT: %v", parent)

	// If there are no errors return the Parent resource
	res := new(bytes.Buffer)
	err = jsonapi.MarshalPayload(res, parent)
	if err != nil {
		log.Println("Problem marshalling the Parent in Show()")
		return c.Render(http.StatusInternalServerError, r.JSON(err.Error()))
	}

	return c.Render(200, r.Func("application/json",
		customJSONRenderer(res.String())))
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
