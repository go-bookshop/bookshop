package main

import (
	"bookshop/internal/httputil"
	"bookshop/internal/models"
	"bookshop/internal/repository"
	"bookshop/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createBooksCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := httputil.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &models.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	v := validator.New()
	if models.ValidateCategory(v, category); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	err = app.repositories.CategoryRepository.Insert(category)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateItem):
			v.AddError("name", fmt.Sprintf("category with the name %q already exists", category.Name))
			app.validationErrorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/categories/%d", category.ID))

	err = httputil.WriteJSON(w, http.StatusCreated, category, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
