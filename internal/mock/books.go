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

func (r *BookRepository) GetBooks(pd *httputil.PaginationData) ([]data.Book, int, error) {
	books := []data.Book{
		{
			ID:        1,
			Title:     "test",
			AvgReview: 3.3,
			ImageUrls: []string{"https://image1.jpg"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "test2",
			AvgReview: 3.4,
			ImageUrls: []string{"https://image2.jpg"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        3,
			Title:     "test3",
			AvgReview: 3.6,
			ImageUrls: []string{"https://image3.jpg"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	booksLen := len(books)

	start := pd.PageNumber * pd.PageSize
	if start >= booksLen {
		return []data.Book{}, 1, nil
	}

	end := start + pd.PageSize
	if end > booksLen {
		end = booksLen
	}

	return books[start:end], (booksLen + pd.PageSize - 1) / pd.PageSize, nil
}
