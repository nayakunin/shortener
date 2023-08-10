package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestFileStorage(t *testing.T) {
	filePath := "test_storage.txt"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			// do nothing
			fmt.Println(err)
		}
	}(filePath) // clean up the test file after the test

	t.Run("newFileStorage", func(t *testing.T) {
		res := newFileStorage(filePath)
		assert.Equal(t, &FileStorage{
			fileStoragePath: filePath,
			Storage: Storage{
				links: make(map[string]interfaces.Link),
				users: make(map[string][]interfaces.Link),
			},
		}, &res)
	})

	t.Run("Add and Get", func(t *testing.T) {
		storage := newFileStorage(filePath)
		err := storage.restoreData()
		assert.NoError(t, err)

		key, err := storage.Add("https://www.example.com", "user1")
		assert.NoError(t, err)

		retrievedLink, err := storage.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, "https://www.example.com", retrievedLink)
	})

	t.Run("AddBatch", func(t *testing.T) {
		storage := newFileStorage(filePath)
		err := storage.restoreData()
		assert.NoError(t, err)

		batches := []interfaces.BatchInput{
			{OriginalURL: "https://www.example1.com"},
			{OriginalURL: "https://www.example2.com"},
		}

		output, err := storage.AddBatch(batches, "user1")
		assert.NoError(t, err)
		assert.Len(t, output, 2)
	})

	t.Run("DeleteUserUrls", func(t *testing.T) {
		storage := newFileStorage(filePath)
		err := storage.restoreData()
		assert.NoError(t, err)

		key, err := storage.Add("https://www.example.com", "user1")
		assert.NoError(t, err)

		err = storage.DeleteUserUrls("user1", []string{key})
		assert.NoError(t, err)

		_, err = storage.Get(key)
		assert.Error(t, err)
	})
}
