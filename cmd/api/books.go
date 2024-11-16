package main

import (
	"bookshop/internal/data"
	"bookshop/internal/httputil"
	"log/slog"
	"net/http"
)

func (app *application) getBooksHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		PageNumber int             `json:"page"`
		PageSize   int             `json:"size"`
		MaxPages   int             `json:"max_pages"`
		Data       []data.BookItem `json:"data"`
	}

	pagination, err := httputil.ParsePaginationQuery(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	books, booksCount, err := app.repositories.BookRepository.GetBooks(pagination)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("number of all items: ", slog.Int("books_count", booksCount))

	res := response{
		PageNumber: pagination.PageNumber,
		PageSize:   pagination.PageSize,
		MaxPages:   httputil.CalculateMaxPages(booksCount, pagination.PageSize),
		Data:       books,
	}

	err = httputil.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
