package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	result := Encode("https://google.com")

	assert.Equal(t, len(result), 8)
}
