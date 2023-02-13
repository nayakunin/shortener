package utils

import (
	"math"
	"strings"
)

const base = 62
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

func Decode(str string) int64 {
	var result int64
	exponent := float64(len(str) - 1)
	for _, c := range str {
		result += int64(math.Pow(base, exponent)) * int64(strings.IndexByte(charset, byte(c)))
		exponent -= 1
	}
	return result
}
