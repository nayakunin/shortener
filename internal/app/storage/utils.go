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

	links, users := parseCSV(csvData)

	return links, users, nil
}

func parseCSV(csvData [][]string) (map[string]Link, map[string][]Link) {
	links := make(map[string]Link)
	users := make(map[string][]Link)

	for _, row := range csvData {
		if len(row) != 2 {
			log.Printf("invalid row: %v", row)
			continue
		}

		link := Link{
			ShortUrl: row[0],
			LongUrl:  row[1],
			UserId:   row[2],
		}

		links[row[0]] = link
		users[row[2]] = append(users[row[2]], link)

	}

	return links, users
}

func writeLinkToFile(fileStoragePath string, key string, link string, userId string) error {
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

	if err := writer.Write([]string{key, link, userId}); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	writer.Flush()

	return nil
}
