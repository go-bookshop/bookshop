BEGIN;

CREATE TABLE IF NOT EXISTS book_reviews(
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stars SMALLINT NOT NULL CHECK (stars BETWEEN 0 AND 5),
    comment TEXT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (book_id, user_id)
);

CREATE INDEX idx_reviews_book_id ON book_reviews(book_id);
CREATE INDEX idx_reviews_user_id ON book_reviews(user_id);

COMMIT;