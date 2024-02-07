// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package sqlite

import (
	"context"
)

const getPost = `-- name: GetPost :one
SELECT id, shown FROM posts
WHERE id = ? LIMIT 1
`

func (q *Queries) GetPost(ctx context.Context, id string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i Post
	err := row.Scan(&i.ID, &i.Shown)
	return i, err
}

const writePost = `-- name: WritePost :one
INSERT INTO posts (
  id, shown
) VALUES (
  ?, ?
)
RETURNING id, shown
`

type WritePostParams struct {
	ID    string
	Shown bool
}

func (q *Queries) WritePost(ctx context.Context, arg WritePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, writePost, arg.ID, arg.Shown)
	var i Post
	err := row.Scan(&i.ID, &i.Shown)
	return i, err
}
