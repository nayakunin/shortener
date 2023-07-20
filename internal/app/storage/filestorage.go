package storage

import (
	"fmt"
	"os"

	"github.com/nayakunin/shortener/internal/app/utils"
	"github.com/pkg/errors"
)

type FileStorage struct {
	Storage
	fileStoragePath string
}

func newFileStorage(fileStoragePath string) FileStorage {
	return FileStorage{
		Storage: Storage{
			links: make(map[string]Link),
			users: make(map[string][]Link),
		},
		fileStoragePath: fileStoragePath,
	}
}

// Get returns original url by key
func (s *FileStorage) Get(key string) (string, error) {
	s.Lock()
	defer s.Unlock()

	link, ok := s.links[key]

	if !ok {
		return "", ErrKeyNotFound
	}

	if link.IsDeleted {
		return "", ErrKeyDeleted
	}

	return link.OriginalURL, nil
}

// Add adds new link to storage
func (s *FileStorage) Add(link string, userID string) (string, error) {
	key := utils.Encode(link)

	s.Lock()
	defer s.Unlock()

	if _, ok := s.links[key]; ok {
		return key, ErrKeyExists
	}

	linkObject := Link{
		ShortURL:    key,
		OriginalURL: link,
		UserID:      userID,
	}

	s.links[key] = linkObject
	s.users[userID] = append(s.users[userID], linkObject)
	if err := writeLinkToFile(s.fileStoragePath, key, link, userID); err != nil {
		return "", err
	}

	return key, nil
}

// AddBatch adds new links to storage
func (s *FileStorage) AddBatch(batches []BatchInput, userID string) ([]BatchOutput, error) {
	output := make([]BatchOutput, len(batches))
	for i, linkObject := range batches {
		key, err := s.Add(linkObject.OriginalURL, userID)
		if err != nil && !errors.Is(err, ErrKeyExists) {
			return nil, err
		}
		output[i] = BatchOutput{
			Key:           key,
			CorrelationID: linkObject.CorrelationID,
		}
	}

	return output, nil
}

// GetUrlsByUser returns all user's URLs
func (s *FileStorage) GetUrlsByUser(id string) (map[string]string, error) {
	s.Lock()
	defer s.Unlock()

	links := make(map[string]string)
	for _, link := range s.users[id] {
		links[link.ShortURL] = link.OriginalURL
	}

	return links, nil
}

func (s *FileStorage) restoreData() error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v", err)
			return
		}
	}(file)

	links, users, err := readLinksFromFile(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	s.links = links
	s.users = users

	return nil
}

// DeleteUserUrls deletes user's URLs
func (s *FileStorage) DeleteUserUrls(userID string, keys []string) error {
	s.Lock()
	defer s.Unlock()

	userLinks := s.users[userID]

	for _, key := range keys {
		link := s.links[key]
		if link.UserID != userID {
			continue
		}

		s.links[key] = Link{
			ShortURL:    key,
			OriginalURL: link.OriginalURL,
			UserID:      link.UserID,
			IsDeleted:   true,
		}

		for i, userLink := range userLinks {
			if userLink.ShortURL == key {
				userLinks[i] = Link{
					ShortURL:    key,
					OriginalURL: userLink.OriginalURL,
					UserID:      userLink.UserID,
					IsDeleted:   true,
				}
			}
		}
	}

	s.users[userID] = userLinks

	links := make([]Link, 0, len(s.links))
	for _, link := range s.links {
		links = append(links, link)
	}

	if err := writeLinksToFile(s.fileStoragePath, links); err != nil {
		return err
	}

	return nil
}
