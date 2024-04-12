package handler

import (
	"BookTalkTwo/cmd/web"
	"BookTalkTwo/cmd/web/components"
	db "BookTalkTwo/db/sqlc"
	"BookTalkTwo/models"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	books "google.golang.org/api/books/v1"
	"google.golang.org/api/option"
)

type postBookCommentParams struct {
	BookID  string `form:"volume_id"`
	Comment string `form:"comment"`
}

func (h *Handler) PostBookComment(c echo.Context) error {
	req := new(postBookCommentParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	userId := h.Get(c, "user_id").(int64)

	commentReq := db.CreateBookCommentParams{
		BookID:  req.BookID,
		Content: req.Comment,
		UserID:  userId,
	}

	_, err := h.Store.CreateBookComment(c.Request().Context(), commentReq)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create book comment")
	}

	bookComments, err := h.Store.GetBookComments(c.Request().Context(), req.BookID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get book comments")
	}
	return renderTemplate(c, components.BookComments(req.BookID, bookComments))
}

type getBookCommentsParams struct {
	BookID string `param:"id"`
}

func (h *Handler) GetBookCommentsHandler(c echo.Context) error {
	req := new(getBookCommentsParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	bookComments, err := h.Store.GetBookComments(c.Request().Context(), req.BookID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get book comments")
	}
	return renderTemplate(c, components.DisplayBookComments(req.BookID, bookComments))
}

type searchParams struct {
	Search string `form:"search"`
}

func (h *Handler) BookViewerByIdHandler(c echo.Context) error {
	req := new(BookByIdParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	return renderTemplate(c, web.BookViewer(req.ID))
}

func (h *Handler) SearchHandler(c echo.Context) error {

	req := new(searchParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	fmt.Println("Search: ", req.Search)
	token, ok := c.Get("access_token").(*oauth2.Token)
	if !ok {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	ctx := c.Request().Context()
	booksService, err := books.NewService(ctx, option.WithTokenSource(h.Conf.TokenSource(ctx, token)))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create books service")
	}
	books, err := booksService.Volumes.List(req.Search).MaxResults(25).Do()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to search books"+err.Error())
	}
	var volumes []models.Volume
	for _, b := range books.Items {
		authors := ""
		thumbnail := ""
		if b.VolumeInfo.Authors != nil {
			authors = strings.Join(b.VolumeInfo.Authors, " / ")
		}
		if b.VolumeInfo.ImageLinks != nil {
			thumbnail = b.VolumeInfo.ImageLinks.Thumbnail
		}
		volumes = append(volumes, models.Volume{
			ID:        b.Id,
			Title:     b.VolumeInfo.Title,
			Thumbnail: thumbnail,
			Author:    authors,
		})
	}
	return renderTemplate(c, components.SearchResults(volumes))
}

type BookByIdParams struct {
	ID string `param:"id"`
}

func (h *Handler) BookByIdHandler(c echo.Context) error {
	req := new(BookByIdParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	token, ok := c.Get("access_token").(*oauth2.Token)
	if !ok {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	ctx := c.Request().Context()
	booksService, err := books.NewService(ctx, option.WithTokenSource(h.Conf.TokenSource(ctx, token)))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create books service")
	}

	book, err := booksService.Volumes.Get(req.ID).Do()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get book")
	}

	//get book notes
	//get books comments
	bookComments, err := h.Store.GetBookComments(ctx, req.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get book comments")
	}
	volume := models.Volume{
		ID:          book.Id,
		Title:       book.VolumeInfo.Title,
		Thumbnail:   book.VolumeInfo.ImageLinks.Thumbnail,
		Author:      strings.Join(book.VolumeInfo.Authors, " / "),
		Description: template.HTML(book.VolumeInfo.Description),
		Categories:  book.VolumeInfo.Categories,
		Notes:       []db.Note{},
		Comments:    bookComments,
	}

	return renderTemplate(c, web.Volume(volume))
}
