package storage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

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
