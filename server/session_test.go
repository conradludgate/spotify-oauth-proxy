package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignSession(t *testing.T) {
	config.SessionKey = "helloworld"
	sessionID := "foobarbaz"
	signed := SignSessionID(sessionID)

	sID, ok := ValidSessionID(signed)
	assert.True(t, ok)
	assert.Equal(t, sessionID, sID)
}
