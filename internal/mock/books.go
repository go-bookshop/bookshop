package mock

import (
	"bookshop/internal/models"
	"bookshop/internal/pagination"
	"bookshop/internal/repository"
	"time"

	"github.com/govalues/decimal"
)

func NewBookRepository() repository.BookRepositoryInterface {
	return &BookRepository{}
}

type BookRepository struct {
}

func (r *BookRepository) GetBooks(pd *pagination.MetaData, bf *pagination.BookFilters) ([]models.BookItem, int, error) {
	books := []models.BookItem{
		{
			ID:         1,
			Title:      "test",
			Synopsis:   "something very interesting about this book",
			AvgReview:  3.3,
			ImageUrls:  []string{"https://image1.jpg"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Authors:    "A",
			Categories: "C",
			Property: models.BookProperty{
				ID:        1,
				Format:    models.Paperback,
				Available: models.Available,
				Price:     decimal.Hundred,
			},
		},
		{
			ID:         2,
			Title:      "test2",
			Synopsis:   "something very interesting about this book",
			AvgReview:  3.4,
			ImageUrls:  []string{"https://image2.jpg"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Authors:    "B",
			Categories: "C",
			Property: models.BookProperty{
				ID:        2,
				Format:    models.Paperback,
				Available: models.Available,
				Price:     decimal.Hundred,
			},
		},
		{
			ID:         3,
			Title:      "test3",
			Synopsis:   "something very interesting about this book",
			AvgReview:  3.6,
			ImageUrls:  []string{"https://image3.jpg"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Authors:    "A, B",
			Categories: "C",
			Property: models.BookProperty{
				ID:        3,
				Format:    models.Paperback,
				Available: models.Available,
				Price:     decimal.Hundred,
			},
		},
	}
	booksLen := len(books)

	start := pd.PageNumber * pd.PageSize
	if start >= booksLen {
		return []models.BookItem{}, 1, nil
	}

	end := start + pd.PageSize
	if end > booksLen {
		end = booksLen
	}

	return books[start:end], (booksLen + pd.PageSize - 1) / pd.PageSize, nil
}
