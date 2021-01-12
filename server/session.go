package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const SessionCookie = "SESSION_ID"

func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    url.PathEscape(SignSessionID(sessionID)),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
}

func GetSession(r *http.Request) *Session {
	log.Println(r.Cookies())

	cookie, err := r.Cookie(SessionCookie)
	if err != nil {
		return nil
	}

	signed, err := url.PathUnescape(cookie.Value)
	if err != nil {
		return nil
	}

	sessionID, ok := ValidSessionID(signed)
	if !ok {
		return nil
	}

	session := new(Session)
	db.First(session, Session{ID: sessionID})

	if session.ID == "" {
		return nil
	}
	return session
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
