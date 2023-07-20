// Package utils provides utility functions for the application.
package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// Encode encodes the input string to a random 8-character string.
func Encode(input string) string {
	// Generate a random 6-byte sequence
	randBytes := make([]byte, 6)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	// Append the input string to the random bytes
	bytes := append(randBytes, []byte(input)...)

	// Encode the bytes using base64 encoding
	encoded := base64.RawURLEncoding.EncodeToString(bytes)

	// Truncate the encoded string to 8 characters
	if len(encoded) > 8 {
		encoded = encoded[:8]
	}

	return encoded
}
