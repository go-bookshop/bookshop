package models

import (
	"time"

	"github.com/govalues/decimal"
)

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
	Price          decimal.Decimal  `json:"price"`
	Publisher      string           `json:"publisher"`
	TargetAudience string           `json:"target_audience"`
	Illustrator    string           `json:"illustrator"`
	Illustrations  string           `json:"illustrations"`
	PageNumber     int64            `json:"page_number"`
	Available      BookAvailability `json:"available"`
	PublishedAt    time.Time        `json:"published_at"`
}
