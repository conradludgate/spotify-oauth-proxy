package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

func GetSpotifyID(token *oauth2.Token) (string, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return "", errors.New("Could not complete authorization: could not connect to spotify")
	}
	token.SetAuthHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("Could not complete authorization: invalid response from spotify")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("Could not complete authorization: invalid response from spotify")
	}

	data := new(struct {
		ID string `json:"id"`
	})
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.ID == "" {
		return "", errors.New("Could not complete authorization: invalid response from spotify")
	}

	return data.ID, nil
}
