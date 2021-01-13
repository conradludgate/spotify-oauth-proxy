package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	state := SignSessionID(RandomString())
	c.Redirect(http.StatusSeeOther, OauthClient().AuthCodeURL(state))
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

	sessionID, ok := ValidSessionID(c.Query("state"))
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Could not complete authorization: invalid state"))
		return
	}

	token, err := OauthClient().Exchange(c, code)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, errors.New("Could not complete authorization: invalid code"))
		return
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Could not complete authorization: could not connect to spotify"))
		return
	}
	token.SetAuthHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Could not complete authorization: invalid response from spotify"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Could not complete authorization: invalid response from spotify"))
		return
	}

	data := new(struct {
		ID string `json:"id"`
	})
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.ID == "" {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Could not complete authorization: invalid response from spotify"))
		return
	}

	user := new(User)
	db.FirstOrCreate(user, User{ID: data.ID})

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
