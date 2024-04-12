package models

import (
	db "BookTalkTwo/db/sqlc"
	"html/template"
)

type Volume struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Thumbnail   string        `json:"thumbnail"`
	Author      string        `json:"author"`
	Description template.HTML `json:"description"`
	Categories  []string
	Notes       []db.Note
	Comments    []db.GetBookCommentsRow
}
