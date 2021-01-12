package spotifyapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	Expiry      time.Time `json:"expiry"`
}

func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Add(-10 * time.Second).Before(time.Now())
}

// IsValid reports whether t is non-nil, has an AccessToken, and is not expired.
func (t *Token) IsValid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

// Type returns t.TokenType if non-empty, else "Bearer".
func (t *Token) Type() string {
	if strings.EqualFold(t.TokenType, "bearer") {
		return "Bearer"
	}
	if strings.EqualFold(t.TokenType, "mac") {
		return "MAC"
	}
	if strings.EqualFold(t.TokenType, "basic") {
		return "Basic"
	}
	if t.TokenType != "" {
		return t.TokenType
	}
	return "Bearer"
}

func (t *Token) SetAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", t.Type()+" "+t.AccessToken)
}

type API struct {
	UserID string
	APIKey string

	apiURL string
	token  *Token
	rt     http.RoundTripper
	ctx    context.Context
}

func (api *API) RoundTrip(r *http.Request) (*http.Response, error) {
	api.Authenticate(r)
	return api.roundtrip().RoundTrip(r)
}

func (api *API) Authenticate(r *http.Request) error {
	token, err := api.GetToken()
	if err != nil {
		return err
	}
	token.SetAuthHeader(r)
	return nil
}

func (api *API) GetToken() (*Token, error) {
	if !api.token.IsValid() {
		req, err := http.NewRequestWithContext(api.getCtx(), "GET", api.url(), nil)
		if err != nil {
			return nil, fmt.Errorf("GetToken: could not establish connection to server: %w", err)
		}
		req.SetBasicAuth(api.UserID, api.APIKey)

		resp, err := api.roundtrip().RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("GetToken: error making request to api server: %w", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("GetToken: could not read response body: %w", err)
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("GetToken: could not get api token: %s", body)
		}

		api.token = new(Token)
		if err := json.Unmarshal(body, api.token); err != nil {
			return nil, fmt.Errorf("GetToken: invalid response from api: %w", err)
		}

		if !api.token.IsValid() {
			return nil, fmt.Errorf("GetToken: api returned invalid token")
		}
	}
	return api.token, nil
}

func (api *API) WithHTTPTransport(rt http.RoundTripper) *API {
	api.rt = rt
	return api
}

func (api *API) roundtrip() http.RoundTripper {
	if api.rt == nil {
		return http.DefaultTransport
	}
	return api.rt
}

func (api *API) WithURL(url string) *API {
	api.apiURL = url
	return api
}

func (api *API) url() string {
	if api.apiURL == "" {
		return "https://spotify.conradludgate.com/api/token"
	}
	return api.apiURL
}

func (api *API) WithContext(ctx context.Context) *API {
	api.ctx = ctx
	return api
}

func (api *API) getCtx() context.Context {
	if api.ctx == nil {
		return context.Background()
	}
	return api.ctx
}
