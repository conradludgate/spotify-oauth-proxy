package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GetSpotifyID(c *gin.Context, token *oauth2.Token) (string, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		c.Error(err)
		return "", errors.New("Could not complete authorization: could not connect to spotify")
	}
	token.SetAuthHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.Error(err)
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
		c.Error(err)
		return "", errors.New("Could not complete authorization: invalid response from spotify")
	}

	return data.ID, nil
}
