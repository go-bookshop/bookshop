BEGIN;

DROP TABLE IF EXISTS book_authors;

DROP INDEX IF EXISTS idx_authors_book_id;
DROP INDEX IF EXISTS idx_authors_author_id; 

COMMIT;