package utils

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
