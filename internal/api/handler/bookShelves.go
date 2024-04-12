package handler

import (
	"BookTalkTwo/cmd/web"
	"BookTalkTwo/cmd/web/components"
	"BookTalkTwo/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	books "google.golang.org/api/books/v1"
	"google.golang.org/api/option"
)

func (h *Handler) BookshelvesHandler(c echo.Context) error {

	token := c.Get("access_token").(*oauth2.Token)
	ctx := c.Request().Context()
	booksService, err := books.NewService(ctx, option.WithTokenSource(h.Conf.TokenSource(ctx, token)))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create books service")
	}

	var bookshelfs []models.Bookshelf

	// Get the bookshelves
	bookshelves, err := booksService.Mylibrary.Bookshelves.List().Do()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get bookshelves")
	}
	if len(bookshelves.Items) == 0 {
		return c.String(http.StatusNotFound, "No bookshelves found")
	}
	for _, b := range bookshelves.Items {

		/* 		vol, err := booksService.Mylibrary.Bookshelves.Volumes.List(strconv.FormatInt(b.Id, 10)).Do() */
		if err == nil {

			var bookshelf models.Bookshelf
			bookshelf.ID = int(b.Id)
			bookshelf.Title = b.Title
			bookshelf.Access = b.Access
			bookshelf.VolumeCount = int(b.VolumeCount)
			bookshelf.Volumes = make([]models.Volume, 0)
			bookshelfs = append(bookshelfs, bookshelf)

		}
	}

	return renderTemplate(c, web.Bookshelves(bookshelfs))
}

type bookShelfVolumesParams struct {
	BookshelfID int `param:"id"`
	Limit       int `query:"limit"`
}

func (h *Handler) BookshelvesVolumesHandler(c echo.Context) error {

	req := new(bookShelfVolumesParams)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	token := c.Get("access_token").(*oauth2.Token)
	ctx := c.Request().Context()
	booksService, err := books.NewService(ctx, option.WithTokenSource(h.Conf.TokenSource(ctx, token)))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create books service")
	}

	vol, err := booksService.Mylibrary.Bookshelves.Volumes.List(strconv.FormatInt(int64(req.BookshelfID), 10)).MaxResults(int64(req.Limit)).Do()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get volumes")
	}
	var volumes []models.Volume

	for _, v := range vol.Items {
		volumes = append(volumes, models.Volume{
			ID:        v.Id,
			Title:     v.VolumeInfo.Title,
			Thumbnail: v.VolumeInfo.ImageLinks.Thumbnail,
			Author:    strings.Join(v.VolumeInfo.Authors, " / "),
		})
	}

	return renderTemplate(c, components.BookshelfVolumes(volumes))
}
