BEGIN;

DROP TABLE IF EXISTS tokens;

DROP INDEX IF EXISTS idx_reviews_book_id;
DROP INDEX IF EXISTS idx_reviews_user_id;

COMMIT;