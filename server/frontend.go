package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

}

func NewToken(c *gin.Context) {

}
