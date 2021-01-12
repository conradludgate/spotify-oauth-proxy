package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var tmpl *template.Template

func ParseTemplates() {
	var err error
	tmpl, err = template.ParseGlob(filepath.Join(config.FrontendDir, "*.html"))
	if err != nil {
		panic(err)
	}
}

func RegisterFrontend() {
	ParseTemplates()
	http.HandleFunc("/", Home)
	http.HandleFunc("/new", NewTokenPage)
	http.HandleFunc("/token", TokenPage)
}

func Home(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	if user == nil {
		http.Redirect(w, r, "/api/login", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "index.html", user)
}

func NewTokenPage(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r)
	if session == nil {
		http.Redirect(w, r, "/api/login", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "new_token.html", Scopes)
}

func TokenPage(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r)
	if session == nil {
		http.Redirect(w, r, "/api/login", http.StatusSeeOther)
		return
	}
}
