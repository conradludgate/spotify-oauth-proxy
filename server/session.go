package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const SessionCookie = "SESSION_ID"

func SetSessionCookie(c *gin.Context, sessionID string) {
	c.SetCookie(
		SessionCookie,
		url.PathEscape(SignSessionID(sessionID)),
		0,
		"/",
		"",
		true,
		true,
	)
}

func Authenticated(c *gin.Context) {
	user := GetUserFromSession(c)
	if user == nil {
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
		return
	}
	c.Set("user", user)
}

func GetUserFromSession(c *gin.Context) *User {
	cookie, err := c.Cookie(SessionCookie)
	if err != nil {
		return nil
	}

	signed, err := url.PathUnescape(cookie)
	if err != nil {
		return nil
	}

	sessionID, ok := ValidSessionID(signed)
	if !ok {
		return nil
	}

	session := new(Session)
	db.First(session, Session{ID: sessionID})
	if session.UserID == "" {
		return nil
	}

	user := new(User)
	db.First(user, User{ID: session.UserID})
	if user.ID == "" {
		return nil
	}

	return user
}

func SignSessionID(sessionID string) string {
	return sessionID + "|" + base64.RawURLEncoding.EncodeToString(Sign(sessionID))
}
func ValidSessionID(signedSessionID string) (string, bool) {
	i := strings.Index(signedSessionID, "|")
	if i == -1 {
		return "", false
	}

	sessionID := signedSessionID[:i]
	got, err := base64.RawURLEncoding.DecodeString(signedSessionID[i+1:])
	if err != nil {
		return "", false
	}

	return sessionID, hmac.Equal(got, Sign(sessionID))
}
func Sign(s string) []byte {
	mac := hmac.New(sha256.New, []byte(config.SessionKey))
	mac.Write([]byte(s))
	return mac.Sum(nil)
}

// RandomString returns a length 64 base 64 string
func RandomString() string {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
