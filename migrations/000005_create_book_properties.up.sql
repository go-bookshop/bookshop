BEGIN;

CREATE TYPE BOOK_FORMAT AS ENUM('audiobook', 'e-book', 'paperback', 'hardcover');
CREATE TYPE BOOK_AVAILABILITY AS ENUM('yes', 'no', 'upcoming');

CREATE TABLE IF NOT EXISTS book_properties(
    isbn TEXT PRIMARY KEY,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    format BOOK_FORMAT NOT NULL,
    language TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    publisher TEXT NOT NULL,
    target_audience TEXT NOT NULL,
    illustrator TEXT NOT NULL DEFAULT '',
    illustrations TEXT NOT NULL DEFAULT '',
    cover_type TEXT NOT NULL DEFAULT '',
    page_number BIGINT NOT NULL DEFAULT 0,
    available BOOK_AVAILABILITY NOT NULL DEFAULT 'no',
    published_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_book_id on book_properties(book_id);

COMMIT;