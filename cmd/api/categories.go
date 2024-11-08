package main

import (
	"bookshop/internal/data"
	"bookshop/internal/httputil"
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

	category := &data.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	v := validator.New()
	if data.ValidateCategory(v, category); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	err = app.repositories.CategoriesRepository.Insert(category)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateItem):
			v.AddError("name", fmt.Sprintf("category with the name %q already exists", category.Name))
			app.validationErrorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", "/v1/books/categories")

	err = httputil.WriteJSON(w, http.StatusCreated, category, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
