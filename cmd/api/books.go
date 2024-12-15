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

	bp, err := pagination.ParseBookPaginationQuery(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	books, booksCount, err := app.repositories.BookRepository.GetBooks(bp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	res := response{
		PageNumber: bp.PageNumber,
		PageSize:   bp.PageSize,
		MaxPages:   pagination.CalculateMaxPages(booksCount, bp.PageSize),
		Data:       books,
	}

	err = httputil.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
