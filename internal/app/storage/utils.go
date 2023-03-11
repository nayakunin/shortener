package storage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func restoreLinksFromFile(fileStoragePath string) (*FileStorage, error) {
	file, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v", err)
			return
		}
	}(file)

	links, err := readLinksFromFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return &FileStorage{
		fileStoragePath: fileStoragePath,
		links:           links,
	}, nil
}

func readLinksFromFile(file *os.File) (map[string]string, error) {
	reader := csv.NewReader(file)
	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return convertToMap(csvData), nil
}

func convertToMap(csvData [][]string) map[string]string {
	links := make(map[string]string)

	for _, row := range csvData {
		links[row[0]] = row[1]
	}

	return links
}

func writeLinkToFile(fileStoragePath string, key string, link string) error {
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

	if err := writer.Write([]string{key, link}); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	writer.Flush()

	return nil
}
