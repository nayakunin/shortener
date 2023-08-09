package storage

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestParseCSV(t *testing.T) {
	type Want struct {
		links map[string]interfaces.Link
		users map[string][]interfaces.Link
		error error
	}

	tests := []struct {
		name  string
		input [][]string
		want  Want
	}{
		{
			name:  "empty input",
			input: [][]string{},
			want: Want{
				links: map[string]interfaces.Link{},
				users: map[string][]interfaces.Link{},
				error: nil,
			},
		},
		{
			name: "one link",
			input: [][]string{
				{"key", "url", "user", "true"},
			},
			want: Want{
				links: map[string]interfaces.Link{
					"key": {
						ShortURL:    "key",
						OriginalURL: "url",
						UserID:      "user",
						IsDeleted:   true,
					},
				},
				users: map[string][]interfaces.Link{
					"user": {
						{
							ShortURL:    "key",
							OriginalURL: "url",
							UserID:      "user",
							IsDeleted:   true,
						},
					},
				},
				error: nil,
			},
		},
		{
			name: "two links",
			input: [][]string{
				{"key1", "url1", "user1", "true"},
				{"key2", "url2", "user2", "false"},
			},
			want: Want{
				links: map[string]interfaces.Link{
					"key1": {
						ShortURL:    "key1",
						OriginalURL: "url1",
						UserID:      "user1",
						IsDeleted:   true,
					},
					"key2": {
						ShortURL:    "key2",
						OriginalURL: "url2",
						UserID:      "user2",
						IsDeleted:   false,
					},
				},
				users: map[string][]interfaces.Link{
					"user1": {
						{
							ShortURL:    "key1",
							OriginalURL: "url1",
							UserID:      "user1",
							IsDeleted:   true,
						},
					},
					"user2": {
						{
							ShortURL:    "key2",
							OriginalURL: "url2",
							UserID:      "user2",
							IsDeleted:   false,
						},
					},
				},
				error: nil,
			},
		},
		{
			name: "invalid is_deleted",
			input: [][]string{
				{"key", "url", "user", "invalid"},
			},
			want: Want{
				error: ErrBadCSVFormat,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links, users, err := parseCSV(tt.input)
			assert.ErrorIs(t, err, tt.want.error)
			assert.Equal(t, tt.want.links, links)
			assert.Equal(t, tt.want.users, users)
		})
	}
}

func TestReadLinksFromFile(t *testing.T) {
	type Want struct {
		links map[string]interfaces.Link
		users map[string][]interfaces.Link
		error error
	}

	tests := []struct {
		name  string
		input []string
		want  Want
	}{
		{
			name:  "empty input",
			input: []string{},
			want: Want{
				links: map[string]interfaces.Link{},
				users: map[string][]interfaces.Link{},
			},
		},
		{
			name:  "invalid input",
			input: []string{"key", "url", "user"},
			want: Want{
				error: ErrBadCSVFormat,
			},
		},
		{
			name:  "one link",
			input: []string{"key", "url", "user", "true"},
			want: Want{
				links: map[string]interfaces.Link{
					"key": {
						ShortURL:    "key",
						OriginalURL: "url",
						UserID:      "user",
						IsDeleted:   true,
					},
				},
				users: map[string][]interfaces.Link{
					"user": {
						{
							ShortURL:    "key",
							OriginalURL: "url",
							UserID:      "user",
							IsDeleted:   true,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test-file*.txt")
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			csvWriter := csv.NewWriter(tmpFile)
			err = csvWriter.Write(tt.input)
			assert.NoError(t, err)
			csvWriter.Flush()
			err = csvWriter.Error()
			assert.NoError(t, err)

			_, err = tmpFile.Seek(0, 0)
			assert.NoError(t, err)

			links, users, err := readLinksFromFile(tmpFile)
			assert.ErrorIs(t, err, tt.want.error)
			assert.Equal(t, tt.want.links, links)
			assert.Equal(t, tt.want.users, users)
		})
	}
}

func TestWriteLinksToFile(t *testing.T) {
	type Want struct {
		data  string
		error error
	}

	tests := []struct {
		name  string
		input []interfaces.Link
		want  Want
	}{
		{
			name: "empty input",
			want: Want{
				data: "",
			},
		},
		{
			name: "one link",
			input: []interfaces.Link{
				{
					ShortURL:    "key",
					OriginalURL: "url",
					UserID:      "user",
					IsDeleted:   true,
				},
			},
			want: Want{
				data: "key,url,user,true\n",
			},
		},
		{
			name: "two links",
			input: []interfaces.Link{
				{
					ShortURL:    "key1",
					OriginalURL: "url1",
					UserID:      "user1",
					IsDeleted:   true,
				},
				{
					ShortURL:    "key2",
					OriginalURL: "url2",
					UserID:      "user2",
					IsDeleted:   false,
				},
			},
			want: Want{
				data: "key1,url1,user1,true\nkey2,url2,user2,false\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test-file*.txt")
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			err = writeLinksToFile(tmpFile.Name(), tt.input)
			assert.NoError(t, err)

			_, err = tmpFile.Seek(0, 0)
			assert.NoError(t, err)

			data, err := os.ReadFile(tmpFile.Name())
			assert.NoError(t, err)
			assert.Equal(t, tt.want.data, string(data))
		})
	}
}
