package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	oauthSpotify "golang.org/x/oauth2/spotify"
)

func OauthClient() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.ClientRedirectURL,
		Endpoint:     oauthSpotify.Endpoint,
		Scopes:       []string{"user-read-private"},
	}
}

func Login(c *gin.Context) {
	state := NewSignature("login", RandomString())
	c.Redirect(http.StatusSeeOther, OauthClient().AuthCodeURL(state))
}

func PostLogin(c *gin.Context, code, sessionID string) {
	token, err := OauthClient().Exchange(c, code)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Could not complete authorization: invalid code"))
		return
	}

	id, err := GetSpotifyID(token)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user := new(User)
	db.FirstOrCreate(user, User{ID: id})

	db.Delete(Session{
		UserID: user.ID,
	})
	db.Create(&Session{
		ID:     sessionID,
		UserID: user.ID,

		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expires:      token.Expiry,
	})

	SetSessionCookie(c, sessionID)

	c.HTML(http.StatusOK, "redirect.html", gin.H{
		"path": "/home",
		"text": "Complete Login",
	})
}

func CreateToken(c *gin.Context, code, tokenID string) {
	token := new(Token)
	db.Find(token, Token{ID: tokenID})
	if token.ID == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Could not complete authorization: invalid state"))
		return
	}

	oauthToken, err := OauthClient().Exchange(c, code)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Could not complete authorization: invalid code"))
		return
	}

	apiKey := RandomString()
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), config.HashCost)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Could not complete authorization: something went wrong creating the token"))
		return
	}

	token.APIKeyHash = hash
	token.AccessToken = oauthToken.AccessToken
	token.RefreshToken = oauthToken.RefreshToken
	token.Expires = oauthToken.Expiry
	token.TokenType = oauthToken.TokenType
	db.Save(token)

	c.HTML(http.StatusOK, "token.html", gin.H{
		"name": token.Name,
		"id":   token.ID,
		"key":  apiKey,
	})
}

func SpotifyCallback(c *gin.Context) {
	if err := c.Query("error"); err != "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Could not complete authorization: %s", err))
		return
	}

	code := c.Query("code")
	if code == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("Could not complete authorization: no auth code present"))
		return
	}

	name, value, ok := ValidSignature(c.Query("state"))
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Could not complete authorization: invalid state"))
		return
	}

	switch name {
	case "login":
		PostLogin(c, code, value)
	case "token":
		CreateToken(c, code, value)
	default:
		c.AbortWithError(http.StatusUnauthorized, errors.New("Could not complete authorization: invalid state"))
		return
	}
}
