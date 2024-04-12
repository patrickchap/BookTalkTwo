package server

import (
	"BookTalkTwo/cmd/web"
	"BookTalkTwo/internal/api/handler"
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/oauth2"
)

var (
	conf     *oauth2.Config
	handlers *handler.Handler
	verifier string
)

const (
	SessIDKey  = "uuid"
	SessExpKey = "expireAt"
)

func (s *Server) RegisterRoutes() http.Handler {
	conf = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/books"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("OAUTH_AUTH_URL"),
			TokenURL: os.Getenv("OAUTH_TOKEN_URL"),
		},
		RedirectURL: os.Getenv("OAUTH_REDIRECT_URI"),
	}
	verifier = oauth2.GenerateVerifier()

	handlers := &handler.Handler{
		Store:    s.store,
		Conf:     conf,
		Verifier: verifier,
	}
	handlers.New(handler.CookieOpts{
		Name:   "auth-session",
		Secret: "super-secret-key",
		MaxAge: 86400 * 7,
	}, func(c echo.Context) error {
		return c.Redirect(http.StatusFound,
			c.Echo().Reverse("/"))
	})

	e := echo.New()

	e.Static("/static", "static")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/js/*", echo.WrapHandler(fileServer))

	e.GET("/browse", s.browseHandler, s.Authorize)
	e.GET("/bookshelves", handlers.BookshelvesHandler, s.Authorize)
	e.GET("/bookshelves/:id", handlers.BookshelvesVolumesHandler, s.Authorize)
	e.GET("/browse/search", s.BrowseSearchTab, s.Authorize)
	e.GET("/browse/categories", s.BrowseCategoriesTab, s.Authorize)
	e.POST("/search", handlers.SearchHandler, s.Authorize)
	e.GET("/books/:id", handlers.BookByIdHandler, s.Authorize)
	e.GET("/comments/book/:id", handlers.GetBookCommentsHandler, s.Authorize)
	e.GET("/books/view/:id", handlers.BookViewerByIdHandler, s.Authorize)
	e.POST("/books/comment", handlers.PostBookComment, s.Authorize)

	e.POST("/login", handlers.LoginHandler)
	e.GET("/login/callback", handlers.LoginCallbackHandler)

	e.GET("/logout", handlers.LogoutHandler)
	e.GET("/", func(c echo.Context) error {
		token := handlers.Get(c, "access_token")
		isLoggedIn := false
		if token != nil {
			accessToken := token.(*oauth2.Token)
			isLoggedIn = accessToken.Valid()

			fmt.Println("Access token from home >>>>>>>: ", accessToken)
		}

		return echo.WrapHandler(templ.Handler(web.Home(isLoggedIn)))(c)
	})

	// Protected route that requires authentication
	e.GET("/protected", s.handleProtected, s.Authorize)
	return e
}

func renderTemplate(c echo.Context, component templ.Component) error {
	return echo.WrapHandler(templ.Handler(component))(c)
}

func (s *Server) BrowseSearchTab(c echo.Context) error {
	return renderTemplate(c, web.SearchTab())
}

func (s *Server) BrowseCategoriesTab(c echo.Context) error {
	return renderTemplate(c, web.CategoriesTab())
}

func (s *Server) browseHandler(c echo.Context) error {

	return renderTemplate(c, web.Browse())
}

func (s *Server) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if access token is available in the session or any other storage mechanism
		token := handlers.Get(c, "access_token")
		if token == nil {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		accessToken := token.(*oauth2.Token)
		if accessToken == nil || accessToken.Valid() == false {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		// Set the authenticated client in the Echo context for use in the route handler
		c.Set("access_token", accessToken)

		return next(c)
	}
}

func (s *Server) handleProtected(c echo.Context) error {

	accessToken := handlers.Get(c, "access_token").(*oauth2.Token)
	fmt.Println("Access token from protected >>>>>>>: ", accessToken)
	return c.String(http.StatusOK, "Protected Route")
}

func getAccessTokenFromStorage(c echo.Context) (*oauth2.Token, error) {
	/* sess, _ := session.Get("session", c)
	accessToken := sess.Values["access_token"]
	fmt.Println("Access token from storage >>>>>>>: ", accessToken)
	if sess.Values["access_token"] != nil {
		return sess.Values["access_token"].(*oauth2.Token), nil
	} */
	return nil, fmt.Errorf("Access token not found")
}
