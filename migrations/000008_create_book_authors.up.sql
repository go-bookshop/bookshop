BEGIN;

CREATE TABLE IF NOT EXISTS book_authors (
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

CREATE INDEX idx_authors_book_id ON book_authors(book_id);
CREATE INDEX idx_authors_author_id ON book_authors(author_id);

COMMIT;