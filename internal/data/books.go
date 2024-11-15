package data

import (
	"bookshop/internal/httputil"
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepositoryInterface interface {
	GetBooks(*httputil.PaginationData) ([]BookItem, int, error)
}

func NewBookRepository(DBPool *pgxpool.Pool) BookRepositoryInterface {
	return &BookRepository{DBPool: DBPool}
}

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Synopsis  string    `json:"synopsis"`
	ImageUrls []string  `json:"images"`
	AvgReview float32   `json:"avg_review"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type BookItem struct {
	Book
	Authors    string         `json:"authors"`
	Properties []BookProperty `json:"properties"`
}

type BookFormat string

const (
	Audiobook BookFormat = "audiobook"
	EBook     BookFormat = "e-book"
	Paperback BookFormat = "paperback"
	Hardcover BookFormat = "hardcover"
)

type BookAvailability string

const (
	Available    BookAvailability = "yes"
	NotAvailable BookAvailability = "no"
	Upcoming     BookAvailability = "upcoming"
)

type BookProperty struct {
	ID             int64            `json:"id"`
	ISBN           string           `json:"isbn"`
	BookID         int64            `json:"book_id"`
	Format         BookFormat       `json:"format"`
	Language       string           `json:"language"`
	Price          float64          `json:"price"`
	Publisher      string           `json:"publisher"`
	TargetAudience string           `json:"target_audience"`
	Illustrator    string           `json:"illustrator"`
	Illustrations  string           `json:"illustrations"`
	PageNumber     int64            `json:"page_number"`
	Available      BookAvailability `json:"available"`
	PublishedAt    time.Time        `json:"published_at"`
}

type BookRepository struct {
	DBPool *pgxpool.Pool
}

func (r *BookRepository) GetBooks(pd *httputil.PaginationData) ([]BookItem, int, error) {
	query := `
WITH 
b_authors AS (
    SELECT 
        b.id AS book_id,
        STRING_AGG(DISTINCT a.name, ', ' ORDER BY a.name) AS authors
    FROM 
        books AS b
    LEFT JOIN 
        book_authors AS ba ON ba.book_id = b.id
    LEFT JOIN 
        authors AS a ON ba.author_id = a.id
    GROUP BY 
        b.id
),
b_properties AS (
    SELECT 
        b.id AS book_id,
        JSON_AGG(JSON_BUILD_OBJECT('id', bp.id, 'format', bp.format, 'price', bp.price, 'available', bp.available)) AS properties
    FROM 
        books AS b
    LEFT JOIN 
        book_properties AS bp ON b.id = bp.book_id
    WHERE 
        bp.available = 'yes'
    GROUP BY 
        b.id
)
SELECT 
    b.id, 
    b.title, 
    b.synopsis, 
    b.image_urls, 
    b.avg_review, 
    b.created_at, 
    b.updated_at, 
    ba.authors, 
    bp.properties
FROM 
    books AS b
LEFT JOIN 
    b_authors AS ba ON ba.book_id = b.id
LEFT JOIN 
    b_properties AS bp ON bp.book_id = b.id
WHERE 
	bp.properties IS NOT NULL
OFFSET $1
FETCH FIRST $2 ROWS ONLY;
	`
	countQuery := `
SELECT 
	COUNT(*)
FROM (
	SELECT DISTINCT b.id
	FROM books AS b
	LEFT JOIN book_properties AS bp ON b.id = bp.book_id
	WHERE bp.available = 'yes'
)
	`

	batch := &pgx.Batch{}
	batch.Queue(query, pd.PageNumber*pd.PageSize, pd.PageSize)
	batch.Queue(countQuery)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := r.DBPool.SendBatch(ctx, batch)

	books := []BookItem{}

	rows, err := results.Query()
	if err != nil {
		return nil, -1, err
	}
	defer rows.Close()

	for rows.Next() {
		var b BookItem
		var propertiesJson []byte

		err = rows.Scan(
			&b.ID,
			&b.Title,
			&b.Synopsis,
			&b.ImageUrls,
			&b.AvgReview,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.Authors,
			&propertiesJson,
		)
		if err != nil {
			return nil, -1, err
		}

		if len(propertiesJson) > 0 {
			err = json.Unmarshal(propertiesJson, &b.Properties)
			if err != nil {
				return nil, -1, err
			}
		}

		books = append(books, b)
	}

	var booksCount int
	err = results.QueryRow().Scan(&booksCount)
	if err != nil {
		return nil, -1, err
	}

	if err = results.Close(); err != nil {
		return nil, -1, err
	}

	return books, booksCount, nil
}
