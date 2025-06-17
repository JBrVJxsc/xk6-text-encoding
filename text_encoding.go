package text_encoding

import (
	"encoding/base64"
	"errors"
	"fmt"
	"unicode/utf8"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/text-encoding", new(TextEncoding))
}

// TextEncoding is the main module exposed to k6 JavaScript
type TextEncoding struct{}

// EncodeUTF8 converts a string to UTF-8 bytes
func (TextEncoding) EncodeUTF8(text string) []byte {
	return []byte(text)
}

// EncodeUTF8ToBase64 converts a string to UTF-8 bytes and then to base64
func (TextEncoding) EncodeUTF8ToBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

// DecodeUTF8 converts UTF-8 bytes back to string with validation
func (TextEncoding) DecodeUTF8(data []byte) (string, error) {
	if !utf8.Valid(data) {
		return "", errors.New("invalid UTF-8 bytes")
	}
	return string(data), nil
}

// DecodeUTF8FromBase64 decodes base64 string to UTF-8 text
func (TextEncoding) DecodeUTF8FromBase64(encodedData string) (string, error) {
	if encodedData == "" {
		return "", nil
	}
	decoded, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	if !utf8.Valid(decoded) {
		return "", errors.New("decoded data is not valid UTF-8")
	}
	return string(decoded), nil
}

// CountUTF8Bytes returns the number of bytes in UTF-8 encoding of the string
func (TextEncoding) CountUTF8Bytes(text string) int {
	return len(text)
}

// CountUTF8Runes returns the number of Unicode code points (characters) in the string
func (TextEncoding) CountUTF8Runes(text string) int {
	return utf8.RuneCountInString(text)
}

// IsValidUTF8 checks if the given string is valid UTF-8
func (TextEncoding) IsValidUTF8(text string) bool {
	return utf8.ValidString(text)
}

// IsValidUTF8Bytes checks if the given bytes represent valid UTF-8
func (TextEncoding) IsValidUTF8Bytes(data []byte) bool {
	return utf8.Valid(data)
}
