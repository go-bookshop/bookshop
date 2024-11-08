BEGIN;

DROP TABLE IF EXISTS book_categories;

DROP INDEX IF EXISTS idx_categories_book_id;
DROP INDEX IF  EXISTS idx_categories_category_id;

COMMIT;