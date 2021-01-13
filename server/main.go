package main

import (
	"html/template"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	ConnectDB()
	defer db.Close()

	// RefreshJob()

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"split": func(s string) []string {
			return strings.Split(s, ",")
		},
	})
	r.LoadHTMLGlob(filepath.Join(config.FrontendDir, "*.html"))

	RegisterFrontend(r)
	RegisterAPI(r.Group("/api"))

	r.StaticFile("/", filepath.Join(config.FrontendDir, "index.html"))

	r.Run(":27228")
}
