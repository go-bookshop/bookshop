INSERT INTO books (title, synopsis, image_urls, avg_review)
VALUES
    ('Harry Potter and the Philosophers Stone', 'A young boy discovers he is a wizard and begins his journey at Hogwarts School of Witchcraft and Wizardry.', ARRAY['https://example.com/hp1_cover.jpg', 'https://example.com/hp1_back.jpg'], 4.8),
    ('A Game of Thrones', 'Noble families vie for control of the Iron Throne and the Seven Kingdoms in a land of betrayal and intrigue.', ARRAY['https://example.com/got_cover.jpg'], 4.7),
    ('The Fellowship of the Ring', 'A hobbit and his companions embark on a perilous journey to destroy a powerful ring.', ARRAY['https://example.com/lotr_cover.jpg'], 4.9),
    ('Foundation', 'A mathematician predicts the fall of a galactic empire and works to preserve knowledge.', ARRAY['https://example.com/foundation_cover.jpg'], 4.6),
    ('Murder on the Orient Express', 'A detective investigates a murder aboard a luxurious train.', ARRAY['https://example.com/murder_cover.jpg'], 4.5);

INSERT INTO authors (name, bio)
VALUES
    ('J.K. Rowling', 'British author, best known for the Harry Potter series.'),
    ('George R.R. Martin', 'American novelist and short story writer, best known for A Song of Ice and Fire.'),
    ('J.R.R. Tolkien', 'English writer, poet, and philologist, best known for The Lord of the Rings.'),
    ('Isaac Asimov', 'American writer and professor of biochemistry, best known for science fiction works.'),
    ('Agatha Christie', 'English writer known for her detective novels, particularly those revolving around Hercule Poirot and Miss Marple.');

INSERT INTO book_authors (book_id, author_id)
VALUES
    (1, 1), (2, 2),
    (3, 3), (4, 4),
    (5, 5);

INSERT INTO book_properties (isbn, book_id, format, language, price, publisher, target_audience, illustrator, illustrations, page_number, available, published_at)
VALUES
    ('978-0747532699', 1, 'paperback', 'English', 19.99, 'Bloomsbury', 'Young Adult', '', '', 223, 'yes', '1997-06-26'),
    ('978-0553103540', 2, 'hardcover', 'English', 29.99, 'Bantam Books', 'Adult', '', '', 694, 'yes', '1996-08-06'),
    ('978-0261103573', 3, 'hardcover', 'English', 25.99, 'George Allen & Unwin', 'Adult', '', '', 423, 'yes', '1954-07-29'),
    ('978-0553293357', 4, 'paperback', 'English', 18.99, 'Spectra', 'Adult', '', '', 296, 'yes', '1951-05-01'),
    ('978-0007119318', 5, 'hardcover', 'English', 15.99, 'Collins Crime Club', 'Adult', '', '', 256, 'yes', '1934-01-01');
