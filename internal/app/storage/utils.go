package storage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nayakunin/shortener/internal/app/interfaces"
	"github.com/pkg/errors"
)

// ErrBadCSVFormat is an error that is returned when csv file has bad format
var ErrBadCSVFormat = errors.New("bad csv format")

func readLinksFromFile(file *os.File) (map[string]interfaces.Link, map[string][]interfaces.Link, error) {
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

func parseCSV(csvData [][]string) (map[string]interfaces.Link, map[string][]interfaces.Link, error) {
	links := make(map[string]interfaces.Link)
	users := make(map[string][]interfaces.Link)

	for _, row := range csvData {
		if len(row) != 4 {
			return nil, nil, ErrBadCSVFormat
		}

		isDeleted, err := strconv.ParseBool(row[3])
		if err != nil {
			return nil, nil, ErrBadCSVFormat
		}

		link := interfaces.Link{
			ShortURL:    row[0],
			OriginalURL: row[1],
			UserID:      row[2],
			IsDeleted:   isDeleted,
		}

		links[row[0]] = link
		users[row[2]] = append(users[row[2]], link)

	}

	return links, users, nil
}

func writeLinksToFile(fileStoragePath string, links []interfaces.Link) error {
	file, err := os.OpenFile(fileStoragePath, os.O_TRUNC|os.O_WRONLY, 0644)
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

	for _, link := range links {
		if err := writer.Write([]string{link.ShortURL, link.OriginalURL, link.UserID, strconv.FormatBool(link.IsDeleted)}); err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	writer.Flush()

	return nil
}
