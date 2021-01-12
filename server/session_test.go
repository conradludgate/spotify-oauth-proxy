package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignSession(t *testing.T) {
	sessionID := "foobarbaz"
	signed := SignSessionID(sessionID)

	sID, ok := ValidSessionID(signed)
	assert.Equal(t, sessionID, sID)
	assert.True(t, ok)
}
