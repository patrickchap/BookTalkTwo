// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import (
	"database/sql"
)

type BookComment struct {
	ID        int64        `json:"id"`
	BookID    string       `json:"book_id"`
	UserID    int64        `json:"user_id"`
	Content   string       `json:"content"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type Note struct {
	ID        int64        `json:"id"`
	BookID    string       `json:"book_id"`
	UserID    int64        `json:"user_id"`
	Content   string       `json:"content"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type User struct {
	ID        int64        `json:"id"`
	Username  string       `json:"username"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	FullName  string       `json:"full_name"`
	Email     string       `json:"email"`
	Picture   string       `json:"picture"`
	CreatedAt sql.NullTime `json:"created_at"`
}