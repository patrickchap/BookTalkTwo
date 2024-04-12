
-- name: GetBookById :one
SELECT * FROM book_comments
WHERE book_id = ? LIMIT 1;

-- name: GetBookComments :many
SELECT 
    book_comments.id AS comment_id,
    book_comments.book_id,
    book_comments.content AS comment_content,
    book_comments.created_at AS comment_created_at,
    users.id AS user_id,
    users.username,
    users.first_name,
    users.last_name,
    users.full_name,
    users.email,
    users.picture,
    users.created_at AS user_created_at
FROM 
    book_comments
JOIN 
    users ON book_comments.user_id = users.id
WHERE 
    book_comments.book_id = ?
ORDER BY 
    book_comments.created_at DESC;

-- name: CreateBookComment :one
INSERT INTO book_comments (book_id, user_id, content)
VALUES (?, ?, ?)
RETURNING *;
