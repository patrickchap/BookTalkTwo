// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import (
	"context"
)

type Querier interface {
	CreateBookComment(ctx context.Context, arg CreateBookCommentParams) (BookComment, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetBookById(ctx context.Context, bookID string) (BookComment, error)
	GetBookComments(ctx context.Context, bookID string) ([]GetBookCommentsRow, error)
	GetUser(ctx context.Context, email string) (User, error)
}

var _ Querier = (*Queries)(nil)
