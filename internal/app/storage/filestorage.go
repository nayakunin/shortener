package storage

import (
	"fmt"
	"os"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/nayakunin/shortener/internal/app/utils"
	"github.com/pkg/errors"
)

// FileStorage is a storage that stores data in files
type FileStorage struct {
	Storage
	fileStoragePath string
}

func newFileStorage(fileStoragePath string) FileStorage {
	return FileStorage{
		Storage: Storage{
			links: make(map[string]interfaces.Link),
			users: make(map[string][]interfaces.Link),
		},
		fileStoragePath: fileStoragePath,
	}
}

// Get returns original url by key. Key is a short url
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

	linkObject := interfaces.Link{
		ShortURL:    key,
		OriginalURL: link,
		UserID:      userID,
	}

	s.links[key] = linkObject
	s.users[userID] = append(s.users[userID], linkObject)
	links := []interfaces.Link{{
		ShortURL:    key,
		OriginalURL: link,
		UserID:      userID,
	}}
	if err := writeLinksToFile(s.fileStoragePath, links); err != nil {
		return "", err
	}

	return key, nil
}

// AddBatch adds new links to storage
func (s *FileStorage) AddBatch(batches []interfaces.BatchInput, userID string) ([]interfaces.BatchOutput, error) {
	output := make([]interfaces.BatchOutput, len(batches))
	for i, linkObject := range batches {
		key, err := s.Add(linkObject.OriginalURL, userID)
		if err != nil && !errors.Is(err, ErrKeyExists) {
			return nil, err
		}
		output[i] = interfaces.BatchOutput{
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

		s.links[key] = interfaces.Link{
			ShortURL:    key,
			OriginalURL: link.OriginalURL,
			UserID:      link.UserID,
			IsDeleted:   true,
		}

		for i, userLink := range userLinks {
			if userLink.ShortURL == key {
				userLinks[i] = interfaces.Link{
					ShortURL:    key,
					OriginalURL: userLink.OriginalURL,
					UserID:      userLink.UserID,
					IsDeleted:   true,
				}
			}
		}
	}

	s.users[userID] = userLinks

	links := make([]interfaces.Link, 0, len(s.links))
	for _, link := range s.links {
		links = append(links, link)
	}

	if err := writeLinksToFile(s.fileStoragePath, links); err != nil {
		return err
	}

	return nil
}

// Stats returns stats
func (s *FileStorage) Stats() (interfaces.Stats, error) {
	s.Lock()
	defer s.Unlock()

	return interfaces.Stats{
		Urls:  len(s.links),
		Users: len(s.users),
	}, nil
}
