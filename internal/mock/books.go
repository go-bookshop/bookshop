package mock

import (
	"bookshop/internal/data"
	"bookshop/internal/httputil"
	"time"
)

func NewBookRepository() data.BookRepositoryInterface {
	return &BookRepository{}
}

type BookRepository struct {
}

func (r *BookRepository) GetBooks(pd *httputil.PaginationData) ([]data.BookItem, int, error) {
	books := []data.BookItem{
		{
			Book: data.Book{
				ID:        1,
				Title:     "test",
				Synopsis:  "something very interesting about this book",
				AvgReview: 3.3,
				ImageUrls: []string{"https://image1.jpg"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Authors: "A",
			Property: data.BookProperty{
				ID:        1,
				Format:    data.Paperback,
				Available: data.Available,
				Price:     100,
			},
		},
		{
			Book: data.Book{
				ID:        2,
				Title:     "test2",
				Synopsis:  "something very interesting about this book",
				AvgReview: 3.4,
				ImageUrls: []string{"https://image2.jpg"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Authors: "B",
			Property: data.BookProperty{
				ID:        2,
				Format:    data.Paperback,
				Available: data.Available,
				Price:     100,
			},
		},
		{
			Book: data.Book{
				ID:        3,
				Title:     "test3",
				Synopsis:  "something very interesting about this book",
				AvgReview: 3.6,
				ImageUrls: []string{"https://image3.jpg"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Authors: "A, B",
			Property: data.BookProperty{
				ID:        3,
				Format:    data.Paperback,
				Available: data.Available,
				Price:     100,
			},
		},
	}
	booksLen := len(books)

	start := pd.PageNumber * pd.PageSize
	if start >= booksLen {
		return []data.BookItem{}, 1, nil
	}

	end := start + pd.PageSize
	if end > booksLen {
		end = booksLen
	}

	return books[start:end], (booksLen + pd.PageSize - 1) / pd.PageSize, nil
}
