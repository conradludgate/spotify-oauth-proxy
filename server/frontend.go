package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterFrontend(r gin.IRouter) {
	auth := r.Group("/").Use(Authenticated)
	auth.GET("/home", Home)
	auth.GET("/token/new", NewTokenPage)
	auth.GET("/token/id/:id", TokenPage)
	auth.POST("/token/id/:id/revoke", RefreshTokenAPIKey)
	auth.POST("/token/id/:id/delete", DeleteToken)
	auth.POST("/token/", NewToken)
}

func Home(c *gin.Context) {
	user := c.MustGet("user").(*User)
	db.Model(user).Related(&user.Tokens)
	c.HTML(http.StatusOK, "home.html", user)
}

func NewTokenPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new_token.html", Scopes)
}

func TokenPage(c *gin.Context) {
	user := c.MustGet("user").(*User)
	token := new(Token)
	db.Find(token, Token{ID: c.Param("id"), UserID: user.ID})
	if token.ID == "" {
		c.String(http.StatusNotFound, "token not found")
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "token.html", token)
}

func NewToken(c *gin.Context) {
	name, ok := c.GetPostForm("name")
	if !ok {
		c.String(http.StatusBadRequest, "token must have a name")
		c.Abort()
	}
	scopes, ok := c.GetPostFormArray("scopes")
	if !ok {
		c.String(http.StatusBadRequest, "token must have a list of scopes")
		c.Abort()
	}

	user := c.MustGet("user").(*User)

	token := new(Token)
	db.Find(token, Token{
		UserID: user.ID,
		Name:   name,
	})

	// If the token has an api key, then it's already taken
	if len(token.APIKeyHash) > 0 {
		c.String(http.StatusBadRequest, "name already taken")
		c.Abort()
	}

	id := uuid.New().String()
	if token.ID == "" {
		db.Create(&Token{
			ID:     id,
			Name:   name,
			UserID: user.ID,
			Scopes: scopes,
		})
	} else {
		id = token.ID
	}

	state := NewSignature("token", id)
	c.Redirect(http.StatusSeeOther, OauthClient(scopes...).AuthCodeURL(state))
}

func RefreshTokenAPIKey(c *gin.Context) {
	user := c.MustGet("user").(*User)

	token := new(Token)
	db.Find(token, Token{
		UserID: user.ID,
		ID:     c.Param("id"),
	})

	if len(token.APIKeyHash) == 0 {
		c.String(http.StatusUnauthorized, "cannot revoke api key")
		c.Abort()
	}

	apiKey := RandomString()
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), config.HashCost)
	if err != nil {
		c.String(http.StatusUnauthorized, "Could not complete authorization: something went wrong creating the token")
		c.Error(err)
		c.Abort()
		return
	}

	token.APIKeyHash = hash
	db.Save(token)

	c.HTML(http.StatusOK, "token.html", gin.H{
		"Name":   token.Name,
		"ID":     token.ID,
		"APIKey": apiKey,
		"Scopes": token.Scopes,
	})
}

func DeleteToken(c *gin.Context) {
	user := c.MustGet("user").(*User)

	db.Delete(Token{
		UserID: user.ID,
		ID:     c.Param("id"),
	})

	c.Redirect(http.StatusSeeOther, "/home")
}
