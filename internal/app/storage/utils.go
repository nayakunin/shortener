package storage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func readLinksFromFile(file *os.File) (map[string]Link, map[string][]Link, error) {
	reader := csv.NewReader(file)
	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	links, users, err := parseCSV(csvData)
	if err != nil {
		return nil, nil, err
	}

	return links, users, nil
}

func parseCSV(csvData [][]string) (map[string]Link, map[string][]Link, error) {
	links := make(map[string]Link)
	users := make(map[string][]Link)

	for _, row := range csvData {
		if len(row) != 3 {
			return nil, nil, fmt.Errorf("invalid row: %v", row)
		}

		link := Link{
			ShortURL:    row[0],
			OriginalURL: row[1],
			UserID:      row[2],
		}

		links[row[0]] = link
		users[row[2]] = append(users[row[2]], link)

	}

	return links, users, nil
}

func writeLinkToFile(fileStoragePath string, key string, link string, userID string) error {
	file, err := os.OpenFile(fileStoragePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("error closing file: %v", err)
			return
		}
	}(file)

	writer := csv.NewWriter(file)

	if err := writer.Write([]string{key, link, userID}); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	writer.Flush()

	return nil
}
