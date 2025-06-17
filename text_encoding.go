package text_encoding

import (
	"encoding/base64"
	"errors"
	"fmt"
	"unicode/utf8"

	"go.k6.io/k6/js/modules"
)

// Error messages
const (
	ErrInvalidUTF8       = "invalid UTF-8 bytes"
	ErrInvalidBase64     = "failed to decode base64"
	ErrInvalidUTF8Base64 = "decoded data is not valid UTF-8"
	ErrEmptyInput        = "empty input"
	ErrNilInput          = "nil input"
)

// MaxInputSize is the maximum size of input strings to prevent memory issues
const MaxInputSize = 100 * 1024 * 1024 // 100MB

func init() {
	modules.Register("k6/x/text-encoding", new(TextEncoding))
}

// TextEncoding is the main module exposed to k6 JavaScript.
// It provides functions for encoding and decoding text in various formats.
type TextEncoding struct{}

// EncodeUTF8 converts a string to UTF-8 bytes.
// It validates the input and returns an error if the input is invalid.
func (TextEncoding) EncodeUTF8(text string) ([]byte, error) {
	if text == "" {
		return []byte{}, nil
	}
	if len(text) > MaxInputSize {
		return nil, fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	if !utf8.ValidString(text) {
		return nil, errors.New(ErrInvalidUTF8)
	}
	return []byte(text), nil
}

// EncodeUTF8ToBase64 converts a string to UTF-8 bytes and then to base64.
// It validates the input and returns an error if the input is invalid.
func (TextEncoding) EncodeUTF8ToBase64(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	if len(text) > MaxInputSize {
		return "", fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	if !utf8.ValidString(text) {
		return "", errors.New(ErrInvalidUTF8)
	}
	return base64.StdEncoding.EncodeToString([]byte(text)), nil
}

// DecodeUTF8 converts UTF-8 bytes back to string with validation.
// It returns an error if the input is invalid UTF-8.
func (TextEncoding) DecodeUTF8(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}
	if len(data) > MaxInputSize {
		return "", fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	if !utf8.Valid(data) {
		return "", errors.New(ErrInvalidUTF8)
	}
	return string(data), nil
}

// DecodeUTF8FromBase64 decodes base64 string to UTF-8 text.
// It validates both the base64 encoding and the resulting UTF-8.
func (TextEncoding) DecodeUTF8FromBase64(encodedData string) (string, error) {
	if encodedData == "" {
		return "", nil
	}
	if len(encodedData) > MaxInputSize {
		return "", fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	decoded, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ErrInvalidBase64, err)
	}
	if !utf8.Valid(decoded) {
		return "", errors.New(ErrInvalidUTF8Base64)
	}
	return string(decoded), nil
}

// CountUTF8Bytes returns the number of bytes in UTF-8 encoding of the string.
// It validates the input and returns an error if the input is invalid.
func (TextEncoding) CountUTF8Bytes(text string) (int, error) {
	if text == "" {
		return 0, nil
	}
	if len(text) > MaxInputSize {
		return 0, fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	if !utf8.ValidString(text) {
		return 0, errors.New(ErrInvalidUTF8)
	}
	return len(text), nil
}

// CountUTF8Runes returns the number of Unicode code points (characters) in the string.
// It validates the input and returns an error if the input is invalid.
func (TextEncoding) CountUTF8Runes(text string) (int, error) {
	if text == "" {
		return 0, nil
	}
	if len(text) > MaxInputSize {
		return 0, fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	if !utf8.ValidString(text) {
		return 0, errors.New(ErrInvalidUTF8)
	}
	return utf8.RuneCountInString(text), nil
}

// IsValidUTF8 checks if the given string is valid UTF-8.
// It returns an error if the input size exceeds the maximum allowed size.
func (TextEncoding) IsValidUTF8(text string) (bool, error) {
	if text == "" {
		return true, nil
	}
	if len(text) > MaxInputSize {
		return false, fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	return utf8.ValidString(text), nil
}

// IsValidUTF8Bytes checks if the given bytes represent valid UTF-8.
// It returns an error if the input size exceeds the maximum allowed size.
func (TextEncoding) IsValidUTF8Bytes(data []byte) (bool, error) {
	if len(data) == 0 {
		return true, nil
	}
	if len(data) > MaxInputSize {
		return false, fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	return utf8.Valid(data), nil
}

// Helper function to validate input size
func validateInputSize(size int) error {
	if size > MaxInputSize {
		return fmt.Errorf("input size exceeds maximum allowed size of %d bytes", MaxInputSize)
	}
	return nil
}

// Helper function to validate UTF-8 string
func validateUTF8String(text string) error {
	if !utf8.ValidString(text) {
		return errors.New(ErrInvalidUTF8)
	}
	return nil
}

// Helper function to validate UTF-8 bytes
func validateUTF8Bytes(data []byte) error {
	if !utf8.Valid(data) {
		return errors.New(ErrInvalidUTF8)
	}
	return nil
}
