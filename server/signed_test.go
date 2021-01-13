package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignSession(t *testing.T) {
	config.SessionKey = "helloworld"
	sessionID := "foobarbaz"
	signed := NewSignature("login", sessionID)

	name, sID, ok := ValidSignature(signed)
	assert.True(t, ok)
	assert.Equal(t, sessionID, sID)
	assert.Equal(t, "login", name)
}
