package main

import (
	"bookshop/internal/httputil"
	"bookshop/internal/models"
	"bookshop/internal/validator"
	"fmt"
	"net/http"
)

func (app *application) createAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}
	err := httputil.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	author := &models.Author{
		Name: input.Name,
		Bio:  input.Bio,
	}

	v := validator.New()
	if models.ValidateAuthor(v, author); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	err = app.repositories.AuthorRepository.Insert(author)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/authors/%d", author.ID))

	err = httputil.WriteJSON(w, http.StatusCreated, author, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
