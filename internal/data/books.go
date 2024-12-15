package data

import (
	"bookshop/internal/models"
	"bookshop/internal/pagination"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepositoryInterface interface {
	GetBooks(*pagination.BookPaginationData) ([]models.BookItem, int, error)
}

func NewBookRepository(DBPool *pgxpool.Pool) BookRepositoryInterface {
	return &BookRepository{DBPool: DBPool}
}

type BookRepository struct {
	DBPool *pgxpool.Pool
}

func (r *BookRepository) GetBooks(pd *pagination.BookPaginationData) ([]models.BookItem, int, error) {
	query := `
SELECT 
    b.id, 
    b.title, 
    b.synopsis, 
    b.image_urls, 
    b.avg_review, 
    b.created_at, 
    b.updated_at,
    STRING_AGG(DISTINCT a.name, ', ' ORDER BY a.name) AS authors,
	bp.id,
	bp.isbn,
	bp.format,
	bp.language,
	bp.price,
	bp.publisher,
	bp.target_audience,
	bp.illustrator,
	bp.illustrations,
	bp.page_number,
	bp.available,
	bp.published_at
FROM 
    books AS b
INNER JOIN
	book_authors AS ba ON ba.book_id = b.id
INNER JOIN
	authors AS a ON a.id = ba.author_id
INNER JOIN
	book_properties AS bp ON bp.book_id = b.id
` + pd.BuildFilterQuery() + `	
GROUP BY
	b.id, bp.id
` + pd.BuildSortingQuery() + `
OFFSET $1
FETCH FIRST $2 ROWS ONLY;
	`
	countQuery := `
SELECT 
	COUNT(*)
FROM 
	books AS b
INNER JOIN
	book_properties AS bp ON bp.book_id = b.id
` + pd.BuildFilterQuery()

	batch := &pgx.Batch{}
	batch.Queue(query, pd.PageNumber*pd.PageSize, pd.PageSize)
	batch.Queue(countQuery)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := r.DBPool.SendBatch(ctx, batch)

	books := []models.BookItem{}

	rows, err := results.Query()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		b := models.BookItem{Property: models.BookProperty{}}

		err = rows.Scan(
			&b.ID,
			&b.Title,
			&b.Synopsis,
			&b.ImageUrls,
			&b.AvgReview,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.Authors,
			&b.Property.ID,
			&b.Property.ISBN,
			&b.Property.Format,
			&b.Property.Language,
			&b.Property.Price,
			&b.Property.Publisher,
			&b.Property.TargetAudience,
			&b.Property.Illustrator,
			&b.Property.Illustrations,
			&b.Property.PageNumber,
			&b.Property.Available,
			&b.Property.PublishedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		books = append(books, b)
	}

	var booksCount int
	err = results.QueryRow().Scan(&booksCount)
	if err != nil {
		return nil, 0, err
	}

	if err = results.Close(); err != nil {
		return nil, 0, err
	}

	return books, booksCount, nil
}
