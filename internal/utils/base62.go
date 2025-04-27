package utils

var base62Chars = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

const base = 62

func EncodeBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var encoded []rune
	for num > 0 {
		remainder := num % base
		encoded = append([]rune{base62Chars[remainder]}, encoded...)
		num = num / base
	}

	// Pad with '0' if shorter than 6 characters
	for len(encoded) < 6 {
		encoded = append([]rune{base62Chars[0]}, encoded...)
	}

	return string(encoded)
}
