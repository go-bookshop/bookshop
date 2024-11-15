INSERT INTO books(title, image_urls, avg_review)
VALUES 
    ('The Art of Coding', ARRAY['https://example.com/images/coding1.jpg', 'https://example.com/images/coding2.jpg'], 4.5),
    ('Data Science Essentials', ARRAY['https://example.com/images/data1.jpg', 'https://example.com/images/data2.jpg'], 4.2),
    ('Mastering Golang', ARRAY['https://example.com/images/go1.jpg', 'https://example.com/images/go2.jpg', 'https://example.com/images/go3.jpg'], 4.8),
    ('Machine Learning Basics', ARRAY['https://example.com/images/ml1.jpg'], 4.3),
    ('Advanced Databases', ARRAY['https://example.com/images/db1.jpg', 'https://example.com/images/db2.jpg'], 4.6),
    ('Introduction to Algorithms', ARRAY['https://example.com/images/algo1.jpg', 'https://example.com/images/algo2.jpg'], 4.1);
