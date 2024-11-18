package data

import (
	"bookshop/internal/httputil"
	"context"
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
	Authors  string       `json:"authors"`
	Property BookProperty `json:"property"`
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

func BookItemSortKeyValidator(key string) bool {
	if key != "avg_review" &&
		key != "created_at" &&
		key != "price" {
		return false
	}

	return true
}

func (r *BookRepository) GetBooks(pd *httputil.PaginationData) ([]BookItem, int, error) {
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
WHERE
	bp.available = $1
GROUP BY
	b.id, bp.id
` + pd.BuildSortingQuery() + `
OFFSET $2
FETCH FIRST $3 ROWS ONLY;
	`
	countQuery := `
SELECT 
	COUNT(*)
FROM 
	books AS b
INNER JOIN
	book_properties AS bp ON bp.book_id = b.id
WHERE bp.available = $1
	`

	batch := &pgx.Batch{}
	batch.Queue(query, Available, pd.PageNumber*pd.PageSize, pd.PageSize)
	batch.Queue(countQuery, Available)

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
		b := BookItem{Property: BookProperty{}}

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
			return nil, -1, err
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
