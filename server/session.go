package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
)

func GetSession(r *http.Request) *Session {
	cookie, err := r.Cookie("SESSION_ID")
	if err != nil {
		return nil
	}

	i := strings.Index(cookie.Value, "|")
	if i == -1 {
		return nil
	}
	sessionID := cookie.Value[:i]
	if !ValidSessionID(sessionID, cookie.Value[i:]) {
		return nil
	}

	session := new(Session)
	db.First(session, "id = ?", sessionID)

	if session.ID == "" {
		return nil
	}
	return session
}

func SignSessionID(sessionID string) string {
	return sessionID + "|" + base64.RawStdEncoding.EncodeToString(Sign(sessionID))
}
func ValidSessionID(sessionID, mac string) bool {
	got, err := base64.RawURLEncoding.DecodeString(mac)
	if err != nil {
		return false
	}

	return hmac.Equal(got, Sign(sessionID))
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
