package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	oauthSpotify "golang.org/x/oauth2/spotify"
)

const id = "SpotifyAuthProxy"

var (
	oauth  *oauth2.Config
	domain string
)

func main() {
	ConnectDB()
	defer db.Close()

	oauth = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.ClientRedirectURL,
		Endpoint:     oauthSpotify.Endpoint,
		Scopes:       []string{"user-read-private"},
	}

	redirect, err := url.Parse(config.ClientRedirectURL)
	if err != nil {
		panic(err)
	}
	domain = redirect.Hostname()

	RegisterFrontend()

	http.HandleFunc("/api/data", Data)
	http.HandleFunc("/api/login", Login)
	http.HandleFunc("/api/spotify_callback", SpotifyCallback)
	http.HandleFunc("/api/new", NewToken)

	log.Fatal(http.ListenAndServe(":27228", nil))
}

func Login(w http.ResponseWriter, r *http.Request) {
	if GetSession(r) != nil {
		http.Redirect(w, r, "/api/data", http.StatusSeeOther)
		return
	}

	state := SignSessionID(RandomString())
	http.Redirect(w, r, oauth.AuthCodeURL(state), http.StatusSeeOther)
}

func SpotifyCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.FormValue("error"); err != "" {
		fmt.Fprintf(w, "Could not complete authorization: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		fmt.Fprintln(w, "Could not complete authorization: no auth code present")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionID, ok := ValidSessionID(r.FormValue("state"))
	if !ok {
		fmt.Fprintln(w, "Could not complete authorization: invalid state")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := oauth.Exchange(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Could not complete authorization: invalid code")
		return
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not complete authorization: could not connect to spotify")
		return
	}
	token.SetAuthHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not complete authorization: invalid response from spotify")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not complete authorization: invalid response from spotify")
		return
	}

	data := new(struct {
		ID string `json:"id"`
	})
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.ID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not complete authorization: invalid response from spotify")
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

	SetSessionCookie(w, sessionID)

	// http.Redirect(w, r, "/api/data", http.StatusSeeOther)
	fmt.Fprintln(w, `<!DOCTYPE html>
<html>
<body>
<script>window.location.replace("/api/data");</script>
<noscript><a href="/api/data">Complete login</a></noscript>
</body>
</hmtl>`)
}

func Data(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	if user == nil {
		http.Redirect(w, r, "/api/login", http.StatusSeeOther)
		return
	}

	fmt.Fprintln(w, user.ID)
}

func NewToken(w http.ResponseWriter, r *http.Request) {

}
