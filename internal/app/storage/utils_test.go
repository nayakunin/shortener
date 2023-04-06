package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCSV(t *testing.T) {
	type Want struct {
		links map[string]Link
		users map[string][]Link
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
				links: map[string]Link{},
				users: map[string][]Link{},
				error: nil,
			},
		},
		{
			name: "one link",
			input: [][]string{
				{"key", "url", "user", "true"},
			},
			want: Want{
				links: map[string]Link{
					"key": {
						ShortURL:    "key",
						OriginalURL: "url",
						UserID:      "user",
						IsDeleted:   true,
					},
				},
				users: map[string][]Link{
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
				links: map[string]Link{
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
				users: map[string][]Link{
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
