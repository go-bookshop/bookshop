CREATE TABLE IF NOT EXISTS book_categories (
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, category_id)
);
