package main

import (
	"bookshop/internal/httputil"
	"bookshop/internal/models"
	"bookshop/internal/pagination"
	"net/http"
)

func (app *application) getBooksHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		PageNumber int               `json:"page"`
		PageSize   int               `json:"size"`
		MaxPages   int               `json:"max_pages"`
		Data       []models.BookItem `json:"data"`
	}

	metadata, filters, err := pagination.ParseBookPaginationQuery(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	books, booksCount, err := app.repositories.BookRepository.GetBooks(metadata, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	res := response{
		PageNumber: metadata.PageNumber,
		PageSize:   metadata.PageSize,
		MaxPages:   pagination.CalculateMaxPages(booksCount, metadata.PageSize),
		Data:       books,
	}

	err = httputil.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
