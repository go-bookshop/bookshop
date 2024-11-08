BEGIN;

DROP TABLE IF EXISTS book_reviews;

DROP INDEX IF EXISTS idx_reviews_book_id;
DROP INDEX IF EXISTS idx_reviews_user_id;

COMMIT;