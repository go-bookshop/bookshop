package main

import (
	"bookshop/internal/data"
	"bookshop/internal/utils"
	"bookshop/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}
	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	author := &data.Author{
		Name: input.Name,
		Bio:  input.Bio,
	}

	v := validator.New()
	if data.ValidateAuthor(v, author); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	err = app.models.AuthorModel.Insert(author)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateAuthorName):
			v.AddError("name", fmt.Sprintf("author with the name %q already exists", author.Name))
			app.validationErrorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/authors/%d", author.ID))

	err = utils.WriteJSON(w, http.StatusCreated, author, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
