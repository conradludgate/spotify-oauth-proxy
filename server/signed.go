package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

func NewSignature(name, value string) string {
	return name + "|" + value + "|" + base64.RawURLEncoding.EncodeToString(Sign(value))
}

func ValidSignature(signed string) (name, value string, ok bool) {
	sections := strings.Split(signed, "|")
	if len(sections) != 3 {
		return
	}

	name = sections[0]
	value = sections[1]
	signature := sections[2]

	got, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return
	}

	ok = hmac.Equal(got, Sign(value))

	return
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
