package main

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	ConnectDB()
	defer db.Close()

	// RefreshJob()

	r := gin.Default()
	r.LoadHTMLGlob(filepath.Join(config.FrontendDir, "*.html"))

	RegisterFrontend(r)
	RegisterAPI(r.Group("/api"))

	r.StaticFile("/", filepath.Join(config.FrontendDir, "index.html"))

	r.Run(":27228")
}
