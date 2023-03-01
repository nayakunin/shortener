package utils

import (
	"encoding/csv"
	"os"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(input string) string {
	inputBytes := []byte(input)
	var num int64
	for i, b := range inputBytes {
		num += int64(b) << uint64(8*i)
	}
	return encodeInt(num)
}

func encodeInt(num int64) string {
	if num == 0 {
		return string(charset[0])
	}
	var result []byte
	chars := []byte(charset)
	length := len(chars)
	for num > 0 {
		result = append(result, chars[num%int64(length)])
		num = num / int64(length)
	}
	return string(result)
}

func ReadLinksFromFile(file *os.File) (map[string]string, error) {
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

func WriteLinkToFile(file *os.File, key string, link string) error {
	writer := csv.NewWriter(file)

	if err := writer.Write([]string{key, link}); err != nil {
		return err
	}

	writer.Flush()

	return nil
}
