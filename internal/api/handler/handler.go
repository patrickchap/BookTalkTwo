package handler

import (
	db "BookTalkTwo/db/sqlc"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type Handler struct {
	Store    db.Store
	Conf     *oauth2.Config
	Verifier string
}

var (
	store *sessions.CookieStore
	mgr   *manager
)

type manager struct {
	session    *sessions.Session
	cookie     CookieOpts
	authFailed echo.HandlerFunc
}

type CookieOpts struct {
	Name   string
	Secret string
	MaxAge int
}

func (h *Handler) New(co CookieOpts, authFailed echo.HandlerFunc) {

	store = sessions.NewCookieStore([]byte(co.Secret))
	sess := sessions.NewSession(store, co.Name)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   co.MaxAge,
		HttpOnly: true,
		Secure:   true,
	}
	mgr = &manager{
		session:    sess,
		cookie:     co,
		authFailed: authFailed,
	}
}

func (h *Handler) GetSession(c echo.Context) *sessions.Session {
	var sess = mgr.session

	if s, err := session.Get(mgr.cookie.Name, c); err == nil {
		sess = s
	}
	return sess
}

func (h *Handler) Get(c echo.Context, key string) any {
	sess := h.GetSession(c)
	return sess.Values[key]
}

func (h *Handler) Set(c echo.Context, key string, value any) error {
	sess := h.GetSession(c)
	sess.Values[key] = value

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Delete(c echo.Context, key string) error {
	sess := h.GetSession(c)
	sess.Options.MaxAge = -1
	delete(sess.Values, key)
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return nil
}
