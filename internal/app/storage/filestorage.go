package storage

import (
	"fmt"
	"os"

	"github.com/nayakunin/shortener/internal/app/utils"
)

type FileStorage struct {
	Storage
	fileStoragePath string
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

func (s *FileStorage) Get(key string) (string, bool) {
	s.Lock()
	defer s.Unlock()

	link, ok := s.links[key]
	return link.ShortUrl, ok
}

func (s *FileStorage) Add(link string, userId string) (string, error) {
	key := utils.Encode(link)

	s.Lock()
	defer s.Unlock()

	if _, ok := s.links[key]; ok {
		return "", ErrKeyExists
	}

	linkObject := Link{
		ShortUrl: key,
		LongUrl:  link,
		UserId:   userId,
	}

	s.links[key] = linkObject
	s.users[userId] = append(s.users[userId], linkObject)
	if err := writeLinkToFile(s.fileStoragePath, key, link, userId); err != nil {
		return "", err
	}

	return key, nil
}

func (s *FileStorage) GetUrlsByUser(id string) (map[string]string, error) {
	s.Lock()
	defer s.Unlock()

	links := make(map[string]string)
	for _, link := range s.users[id] {
		links[link.ShortUrl] = link.LongUrl
	}

	return links, nil
}
