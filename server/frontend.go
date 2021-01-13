package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	oauthSpotify "golang.org/x/oauth2/spotify"
)

func RegisterFrontend(r gin.IRouter) {
	auth := r.Group("/").Use(Authenticated)
	auth.GET("/home", Home)
	auth.GET("/token/new", NewTokenPage)
	auth.GET("/token/id/:id", TokenPage)
	auth.POST("/token/", NewToken)
}

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", c.MustGet("user"))
}

func NewTokenPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new_token.html", Scopes)
}

func TokenPage(c *gin.Context) {
	user := c.MustGet("user").(*User)
	token := new(Token)
	db.Find(token, Token{ID: c.Param("id"), UserID: user.ID})
	if token.ID == "" {
		c.AbortWithError(http.StatusNotFound, errors.New("token now found"))
		return
	}

	c.HTML(http.StatusOK, "token.html", gin.H{
		"name": token.Name,
		"id":   token.ID,
	})
}

func NewToken(c *gin.Context) {
	name, ok := c.GetPostForm("name")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("token must have a name"))
	}
	scopes, ok := c.GetPostFormArray("scopes")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("token must have a list of scopes"))
	}

	oauth := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.ClientRedirectURL,
		Endpoint:     oauthSpotify.Endpoint,
		Scopes:       scopes,
	}

	user := c.MustGet("user").(*User)

	token := new(Token)
	db.Find(token, Token{
		UserID: user.ID,
		Name:   name,
	})

	// If the token has an api key, then it's already taken
	if len(token.APIKeyHash) > 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("name already taken"))
	}

	id := uuid.New().String()
	if token.ID == "" {
		db.Create(&Token{
			ID:     id,
			Name:   name,
			UserID: user.ID,
		})
	} else {
		id = token.ID
	}

	state := NewSignature(name, id)
	c.Redirect(http.StatusSeeOther, oauth.AuthCodeURL(state))
}
