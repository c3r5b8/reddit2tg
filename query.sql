-- name: GetPost :one
SELECT * FROM posts
WHERE id = ? LIMIT 1;

-- name: WritePost :one
INSERT INTO posts (
  id, shown
) VALUES (
  ?, ?
)
RETURNING *;
