package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	// Test when PORT is not set
	os.Unsetenv("PORT")
	assert.NotPanics(t, func() {
		go main()
	})

	// Test when PORT is set to "8081"
	os.Setenv("PORT", "8081")
	assert.NotPanics(t, func() {
		go main()
	})
}
