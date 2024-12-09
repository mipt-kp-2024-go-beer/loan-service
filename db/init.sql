DROP TABLE IF EXISTS lent_books;

CREATE TABLE lent_books (
    id TEXT,
    user_id TEXT,
    book_id TEXT,
    taken_at INTEGER,
    return_deadline INTEGER,
    returned BOOLEAN,
    returned_at INTEGER
);
