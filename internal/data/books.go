package data

import (
	"bookshop/internal/httputil"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepositoryInterface interface {
	GetBooks(*httputil.PaginationData) ([]Book, int, error)
}

func NewBookRepository(DBPool *pgxpool.Pool) BookRepositoryInterface {
	return &BookRepository{DBPool: DBPool}
}

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	ImageUrls []string  `json:"images"`
	AvgReview float32   `json:"avg_review"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type BookRepository struct {
	DBPool *pgxpool.Pool
}

func (r *BookRepository) GetBooks(pd *httputil.PaginationData) ([]Book, int, error) {
	query := `
select b.id, b.title, b.image_urls, b.avg_review, b.created_at, b.updated_at
from books as b
offset $1
fetch first $2 rows only
	`
	countQuery := `
select count(id) from books
	`

	batch := &pgx.Batch{}
	batch.Queue(query, pd.PageNumber*pd.PageSize, pd.PageSize)
	batch.Queue(countQuery)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := r.DBPool.SendBatch(ctx, batch)

	books := []Book{}

	rows, err := results.Query()
	if err != nil {
		return nil, -1, err
	}
	defer rows.Close()

	for rows.Next() {
		var b Book

		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.ImageUrls,
			&b.AvgReview,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, -1, err
		}

		books = append(books, b)
	}

	var booksCount int
	err = results.QueryRow().Scan(&booksCount)
	if err != nil {
		return nil, -1, err
	}

	if err := results.Close(); err != nil {
		return nil, -1, err
	}

	return books, booksCount, nil
}
