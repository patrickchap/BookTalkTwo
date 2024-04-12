package handler

import (
	db "BookTalkTwo/db/sqlc"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"net/http"
	"os"
)

type LoginParams struct {
	Credential string `form:"credential"`
	GCSRFToken string `form:"g_csrf_token" json:"g_csrf_token"`
}

type GoogleClaims struct {
	Email      string `form:"email"`
	Name       string `form:"name"`
	GivenName  string `form:"given_name" json:"given_name"`
	FamilyName string `form:"family_name" json:"family_name"`
	Picture    string `form:"picture"`
	jwt.RegisteredClaims
}

func (h *Handler) LoginHandler(c echo.Context) error {

	resp := new(LoginParams)
	err := c.Bind(resp)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	payload, err := idtoken.Validate(c.Request().Context(), resp.Credential, os.Getenv("OAUTH_CLIENT_ID"))
	if err != nil {
		fmt.Println("Failed to validate token: ", err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	name := payload.Claims["name"]
	email := payload.Claims["email"]
	picture := payload.Claims["picture"]
	givenName := payload.Claims["given_name"]
	familyName := payload.Claims["family_name"]

	args := db.CreateUserParams{
		Username:  email.(string),
		FirstName: givenName.(string),
		LastName:  familyName.(string),
		FullName:  name.(string),
		Email:     email.(string),
		Picture:   picture.(string),
	}

	user, err := h.Store.GetUser(c.Request().Context(), email.(string))
	if err != nil {
		user, err = h.Store.CreateUser(c.Request().Context(), args)
		if err != nil {
			fmt.Println("Failed to create user: ", err)
			return c.String(http.StatusInternalServerError, "Failed to create user")
		}
	}

	// Set the authenticated client in the Echo context
	h.Set(c, "user_id", user.ID)
	h.Set(c, "user_email", email.(string))
	url := h.Conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(h.Verifier))
	return c.Redirect(http.StatusFound, url)
}

func (h *Handler) LogoutHandler(c echo.Context) error {

	h.Set(c, "access_token", nil)
	h.Set(c, "user_email", nil)
	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *Handler) LoginCallbackHandler(c echo.Context) error {
	ctx := c.Request().Context()

	code := c.QueryParam("code")
	if code == "" {
		return c.String(http.StatusBadRequest, "Failed to read code")
	}

	tok, err := h.Conf.Exchange(ctx, code, oauth2.VerifierOption(h.Verifier))
	if err != nil {

		fmt.Println("Failed to exchange token: ", err)
	}

	h.Set(c, "access_token", tok)

	//redirect to /home
	return c.Redirect(http.StatusSeeOther, "/")
}
