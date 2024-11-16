INSERT INTO users (first_name, last_name, email, password, activated)
VALUES ('Alice', 'Smith', 'alice@example.com',
        '\x2432612431302470514f77514d57314b4c334f4e616d646133546d377037315a7a3662745a54786b', TRUE),
       ('Bob', 'Brown', 'bob@example.com',
        '\x2432612431302470514f77514d57314b4c334f4e616d646133546d377037315a7a3662745a54786b', FALSE),
       ('Charlie', 'Davis', 'charlie@example.com',
        '\x2432612431302470514f77514d57314b4c334f4e616d646133546d377037315a7a3662745a54786b', TRUE),
       ('Diana', 'Evans', 'diana@example.com',
        '\x2432612431302470514f77514d57314b4c334f4e616d646133546d377037315a7a3662745a54786b', FALSE);

UPDATE users
SET updated_at = '2024-01-01 12:00:00+00'
WHERE email = 'alice@example.com';

UPDATE users
SET updated_at = '2024-01-01 12:00:00+00'
WHERE email = 'charlie@example.com';
