package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func renderTemplate(c echo.Context, component templ.Component) error {
	return echo.WrapHandler(templ.Handler(component))(c)
}
