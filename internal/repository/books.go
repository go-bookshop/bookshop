package repository

import (
	"bookshop/internal/models"
	"bookshop/internal/pagination"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepositoryInterface interface {
	GetBooks(*pagination.MetaData, *pagination.BookFilters) ([]models.BookItem, int, error)
}

func NewBookRepository(DBPool *pgxpool.Pool) BookRepositoryInterface {
	return &BookRepository{DBPool: DBPool}
}

type BookRepository struct {
	DBPool *pgxpool.Pool
}

func (r *BookRepository) GetBooks(pd *pagination.MetaData, filters *pagination.BookFilters) ([]models.BookItem, int, error) {
	query := `
WITH filtered_books AS (
    SELECT 
        b.id, 
        b.title, 
        b.synopsis, 
        b.image_urls, 
        b.avg_review, 
        b.created_at, 
        b.updated_at,
        STRING_AGG(DISTINCT a.name, ', ' ORDER BY a.name) AS authors,
        bp.id AS book_property_id,
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
    INNER JOIN 
        book_categories AS bc ON bc.book_id = b.id
    INNER JOIN 
        categories AS c ON c.id = bc.category_id
` + filters.BuildFilterQuery() + `
    GROUP BY 
        b.id, bp.id
),
book_categories_agg AS (
    SELECT
        b.id AS book_id,
        STRING_AGG(DISTINCT c.name, ', ' ORDER BY c.name) AS categories
    FROM 
        books AS b
    INNER JOIN 
        book_categories AS bc ON bc.book_id = b.id
    INNER JOIN 
        categories AS c ON c.id = bc.category_id
    WHERE 
        b.id IN (SELECT id FROM filtered_books)
    GROUP BY 
        b.id
)
SELECT 
    COUNT(*) OVER () AS total_items,
    fb.id, 
    fb.title, 
    fb.synopsis, 
    fb.image_urls, 
    fb.avg_review, 
    fb.created_at, 
    fb.updated_at,
    fb.book_property_id,
    fb.isbn,
    fb.format,
    fb.language,
    fb.price,
    fb.publisher,
    fb.target_audience,
    fb.illustrator,
    fb.illustrations,
    fb.page_number,
    fb.available,
    fb.published_at,
	fb.authors,
    bca.categories
FROM 
    filtered_books AS fb
INNER JOIN 
    book_categories_agg AS bca ON bca.book_id = fb.id
` + pd.BuildSortingQuery() + `	
OFFSET $1
FETCH FIRST $2 ROWS ONLY;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var books []models.BookItem
	var booksCount int
	rows, err := r.DBPool.Query(ctx, query, pd.PageNumber*pd.PageSize, pd.PageSize)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		b := models.BookItem{Property: models.BookProperty{}}

		err = rows.Scan(
			&booksCount,
			&b.ID,
			&b.Title,
			&b.Synopsis,
			&b.ImageUrls,
			&b.AvgReview,
			&b.CreatedAt,
			&b.UpdatedAt,
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
			&b.Authors,
			&b.Categories,
		)
		if err != nil {
			return nil, 0, err
		}

		books = append(books, b)
	}

	return books, booksCount, nil
}
