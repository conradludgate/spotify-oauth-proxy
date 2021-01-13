package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterFrontend(r gin.IRouter) {

	auth := r.Group("/", Authenticated)
	auth.GET("/home", Home)
	auth.GET("/new", NewTokenPage)
	auth.GET("/token/:id", TokenPage)
	auth.POST("/token", NewToken)
}

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", c.MustGet("user"))
}

func NewTokenPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new_token.html", Scopes)
}

func TokenPage(c *gin.Context) {

}

func NewToken(c *gin.Context) {

}
