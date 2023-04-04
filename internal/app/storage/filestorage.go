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

func (s *FileStorage) Get(key string) (string, bool) {
	s.Lock()
	defer s.Unlock()

	link, ok := s.links[key]
	return link.OriginalURL, ok
}

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
