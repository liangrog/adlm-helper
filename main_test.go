package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	isD := "testing/"
	assert.True(t, isDir(isD))

	notD := "testing"
	assert.False(t, isDir(notD))
}
