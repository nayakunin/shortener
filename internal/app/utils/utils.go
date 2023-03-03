package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/nayakunin/shortener/internal/app/server/config"
)

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

func ReadLinksFromFile(file *os.File) (map[string]string, error) {
	reader := csv.NewReader(file)
	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return convertToMap(csvData), nil
}

func convertToMap(csvData [][]string) map[string]string {
	links := make(map[string]string)

	for _, row := range csvData {
		links[row[0]] = row[1]
	}

	return links
}

func WriteLinkToFile(key string, link string) error {
	file, err := os.OpenFile(config.Config.FileStoragePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("error closing file: %v", err)
			panic(err)
		}
	}(file)

	writer := csv.NewWriter(file)

	if err := writer.Write([]string{key, link}); err != nil {
		return err
	}

	writer.Flush()

	return nil
}
