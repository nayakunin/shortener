package storage

import (
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	t.Run("add link", func(t *testing.T) {
		s := newStorage()
		testLink := "https://google.com"

		key, err := s.Add(testLink, "user")
		assert.NoError(t, err)

		link, err := s.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, testLink, link)
	})

	t.Run("link not found", func(t *testing.T) {
		s := newStorage()

		_, err := s.Get("key")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("link deleted", func(t *testing.T) {
		s := newStorage()
		testLink := "https://google.com"
		user := "user"

		key, err := s.Add(testLink, user)
		assert.NoError(t, err)

		err = s.DeleteUserUrls(user, []string{key})
		assert.NoError(t, err)

		_, err = s.Get(key)
		assert.ErrorIs(t, err, ErrKeyDeleted)
	})

	t.Run("add batch", func(t *testing.T) {
		s := newStorage()
		user := "user"
		input := []interfaces.BatchInput{
			{
				OriginalURL:   "https://google.com/1",
				CorrelationID: "1",
			},
			{
				OriginalURL:   "https://google.com/2",
				CorrelationID: "2",
			},
		}

		output, err := s.AddBatch(input, user)
		assert.NoError(t, err)
		assert.Len(t, output, 2)

		for i, o := range output {
			link, err := s.Get(o.Key)
			assert.NoError(t, err)
			assert.Equal(t, input[i].OriginalURL, link)
		}
	})

	t.Run("get urls by user", func(t *testing.T) {
		s := newStorage()
		user := "user"
		input := []interfaces.BatchInput{
			{
				OriginalURL:   "https://google.com/1",
				CorrelationID: "1",
			},
			{
				OriginalURL:   "https://google.com/2",
				CorrelationID: "2",
			},
		}

		output, err := s.AddBatch(input, user)
		assert.NoError(t, err)
		assert.Len(t, output, 2)

		urls, err := s.GetUrlsByUser(user)
		assert.NoError(t, err)
		assert.Len(t, urls, 2)

		for i, o := range output {
			assert.Equal(t, input[i].OriginalURL, urls[o.Key])
		}
	})

	t.Run("don't delete other user urls", func(t *testing.T) {
		s := newStorage()
		user1 := "user1"
		user2 := "user2"
		input := []interfaces.BatchInput{
			{
				OriginalURL:   "https://google.com/1",
				CorrelationID: "1",
			},
			{
				OriginalURL:   "https://google.com/2",
				CorrelationID: "2",
			},
		}

		output, err := s.AddBatch(input, user1)
		assert.NoError(t, err)
		assert.Len(t, output, 2)

		err = s.DeleteUserUrls(user2, []string{output[0].Key})
		assert.NoError(t, err)

		for _, o := range output {
			_, err = s.Get(o.Key)
			assert.NoError(t, err)
		}
	})
}
