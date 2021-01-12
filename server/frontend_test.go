package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTemplates(t *testing.T) {
	config.FrontendDir = "../client"
	ParseTemplates()

	found := make(map[string]bool)

	for _, template := range tmpl.Templates() {
		found[template.Name()] = true
	}

	assert.Equal(t, map[string]bool{
		"index.html":     true,
		"new_token.html": true,
		"token.html":     true,
	}, found)
}
