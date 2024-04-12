-- name: GetUser :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  username,
  first_name,
  last_name,
  full_name,
  email,
  picture
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;
