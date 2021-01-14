package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/oauth2"
)

func RefreshJob() {
	period := 10 * time.Minute
	oauth := OauthClient()

	go func() {
		for {
			t := <-time.After(period)
			log.Println("starting token refresh")

			var tokens int
			db.Where("refresh_token NOT NULL AND expires < ?", t).Count(&tokens)

			for i := 0; i < tokens; i++ {
				token := new(Token)
				db.Where("refresh_token NOT NULL AND expires < ?", t).Take(&token)

				oldToken := token.IntoOauth()
				if oldToken.Valid() {
					continue
				}
				newToken, err := oauth.TokenSource(context.Background(), oldToken).Token()
				if err != nil {
					log.Println("could not refresh token", err.Error())
					continue
				}

				token.AccessToken = newToken.AccessToken
				token.RefreshToken = newToken.RefreshToken
				token.Expires = newToken.Expiry
				token.TokenType = newToken.TokenType
				db.Save(token)
			}
		}
	}()
}

func Refresh(ctx context.Context, t *oauth2.Token) (*oauth2.Token, error) {
	return OauthClient().TokenSource(ctx, t).Token()
}
