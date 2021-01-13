package main

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const id = "SpotifyAuthProxy"

var (
	domain string
)

func main() {
	ConnectDB()
	defer db.Close()

	r := gin.Default()
	r.LoadHTMLGlob(filepath.Join(config.FrontendDir, "*.html"))

	RegisterFrontend(r)

	// http.HandleFunc("/api/data", Data)
	r.GET("/api/login", Login)
	r.GET("/api/spotify_callback", SpotifyCallback)

	r.Run(":27228")
}
