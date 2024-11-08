BEGIN;

CREATE TABLE IF NOT EXISTS book_categories (
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_categories_book_id ON book_categories(book_id);
CREATE INDEX IF NOT EXISTS idx_categories_category_id ON book_categories(category_id);

COMMIT;