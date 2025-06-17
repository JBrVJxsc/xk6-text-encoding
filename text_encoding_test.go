package text_encoding

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestEncodeUTF8(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
		},
		{
			name:     "ascii text",
			input:    "hello",
			expected: []byte{'h', 'e', 'l', 'l', 'o'},
		},
		{
			name:     "unicode text",
			input:    "Hello 🌍",
			expected: []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xF0, 0x9F, 0x8C, 0x8D},
		},
		{
			name:     "chinese characters",
			input:    "你好",
			expected: []byte{0xE4, 0xBD, 0xA0, 0xE5, 0xA5, 0xBD},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.EncodeUTF8(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("EncodeUTF8() length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i, b := range result {
				if b != tt.expected[i] {
					t.Errorf("EncodeUTF8()[%d] = %02x, want %02x", i, b, tt.expected[i])
				}
			}
		})
	}
}

func TestEncodeUTF8ToBase64(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple text",
			input:    "hello",
			expected: base64.StdEncoding.EncodeToString([]byte("hello")),
		},
		{
			name:     "unicode text",
			input:    "Hello 🌍",
			expected: base64.StdEncoding.EncodeToString([]byte("Hello 🌍")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.EncodeUTF8ToBase64(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeUTF8ToBase64() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeUTF8(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		name        string
		input       []byte
		expected    string
		expectError bool
	}{
		{
			name:        "empty bytes",
			input:       []byte{},
			expected:    "",
			expectError: false,
		},
		{
			name:        "valid ascii",
			input:       []byte{'h', 'e', 'l', 'l', 'o'},
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "valid unicode",
			input:       []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xF0, 0x9F, 0x8C, 0x8D},
			expected:    "Hello 🌍",
			expectError: false,
		},
		{
			name:        "invalid utf-8 sequence",
			input:       []byte{0xFF, 0xFE},
			expected:    "",
			expectError: true,
		},
		{
			name:        "incomplete utf-8 sequence",
			input:       []byte{0xF0, 0x9F},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := te.DecodeUTF8(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("DecodeUTF8() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeUTF8() unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("DecodeUTF8() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeUTF8FromBase64(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "valid base64 ascii",
			input:       base64.StdEncoding.EncodeToString([]byte("hello")),
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "valid base64 unicode",
			input:       base64.StdEncoding.EncodeToString([]byte("Hello 🌍")),
			expected:    "Hello 🌍",
			expectError: false,
		},
		{
			name:        "invalid base64",
			input:       "invalid base64!@#",
			expected:    "",
			expectError: true,
		},
		{
			name:        "base64 with invalid utf-8",
			input:       base64.StdEncoding.EncodeToString([]byte{0xFF, 0xFE}),
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := te.DecodeUTF8FromBase64(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("DecodeUTF8FromBase64() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeUTF8FromBase64() unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("DecodeUTF8FromBase64() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountUTF8Bytes(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "ascii text",
			input:    "hello",
			expected: 5,
		},
		{
			name:     "unicode text with emoji",
			input:    "Hello 🌍",
			expected: 10, // "Hello " (6 bytes) + 🌍 (4 bytes) = 10 bytes
		},
		{
			name:     "chinese characters",
			input:    "你好",
			expected: 6, // 3 bytes per character
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.CountUTF8Bytes(tt.input)
			if result != tt.expected {
				t.Errorf("CountUTF8Bytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountUTF8Runes(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "ascii text",
			input:    "hello",
			expected: 5,
		},
		{
			name:     "unicode text with emoji",
			input:    "Hello 🌍",
			expected: 7, // 6 ASCII characters + 1 emoji
		},
		{
			name:     "chinese characters",
			input:    "你好",
			expected: 2, // 2 Chinese characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.CountUTF8Runes(tt.input)
			if result != tt.expected {
				t.Errorf("CountUTF8Runes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsValidUTF8(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "valid ascii",
			input:    "hello",
			expected: true,
		},
		{
			name:     "valid unicode",
			input:    "Hello 🌍",
			expected: true,
		},
		{
			name:     "invalid utf-8",
			input:    string([]byte{0xFF, 0xFE}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.IsValidUTF8(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidUTF8() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsValidUTF8Bytes(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: true,
		},
		{
			name:     "valid ascii",
			input:    []byte{'h', 'e', 'l', 'l', 'o'},
			expected: true,
		},
		{
			name:     "valid unicode",
			input:    []byte{'H', 'e', 'l', 'l', 'o', ' ', 0xF0, 0x9F, 0x8C, 0x8D},
			expected: true,
		},
		{
			name:     "invalid utf-8",
			input:    []byte{0xFF, 0xFE},
			expected: false,
		},
		{
			name:     "incomplete utf-8",
			input:    []byte{0xF0, 0x9F},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.IsValidUTF8Bytes(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidUTF8Bytes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLargeTextRoundtrip(t *testing.T) {
	te := &TextEncoding{}

	// Create a large string with various Unicode characters
	var largeText strings.Builder
	// Add ASCII characters
	for i := 0; i < 5000; i++ {
		largeText.WriteString("Hello ")
	}
	// Add Chinese characters
	for i := 0; i < 2500; i++ {
		largeText.WriteString("你好")
	}
	// Add emojis
	for i := 0; i < 1000; i++ {
		largeText.WriteString("🌍")
	}
	// Add mixed content
	for i := 0; i < 500; i++ {
		largeText.WriteString("café ")
	}
	// Add more diverse content
	for i := 0; i < 250; i++ {
		largeText.WriteString("résumé ")
	}
	// Add Korean characters
	for i := 0; i < 250; i++ {
		largeText.WriteString("안녕하세요 ")
	}
	// Add Arabic text
	for i := 0; i < 250; i++ {
		largeText.WriteString("مرحبا ")
	}

	text := largeText.String()

	// Test UTF-8 encoding and decoding
	encoded := te.EncodeUTF8(text)
	decoded, err := te.DecodeUTF8(encoded)
	if err != nil {
		t.Errorf("DecodeUTF8() error: %v", err)
	}
	if decoded != text {
		t.Error("Large text roundtrip failed")
	}

	// Test Base64 encoding and decoding
	base64Encoded := te.EncodeUTF8ToBase64(text)
	base64Decoded, err := te.DecodeUTF8FromBase64(base64Encoded)
	if err != nil {
		t.Errorf("DecodeUTF8FromBase64() error: %v", err)
	}
	if base64Decoded != text {
		t.Error("Large text base64 roundtrip failed")
	}

	// Test byte and rune counting
	byteCount := te.CountUTF8Bytes(text)
	runeCount := te.CountUTF8Runes(text)

	// Verify byte count matches encoded length
	if byteCount != len(encoded) {
		t.Errorf("Byte count mismatch: CountUTF8Bytes() = %d, actual bytes length = %d", byteCount, len(encoded))
	}

	// Verify rune count is less than byte count (since some runes use multiple bytes)
	if runeCount >= byteCount {
		t.Errorf("Rune count (%d) should be less than byte count (%d)", runeCount, byteCount)
	}

	// Verify rune count matches expected count
	expectedRunes := 5000*6 + // ASCII "Hello " (6 runes each)
		2500*2 + // Chinese "你好" (2 runes each)
		1000*1 + // Emoji "🌍" (1 rune each)
		500*5 + // Mixed "café " (5 runes each)
		250*7 + // French "résumé " (7 runes each)
		250*6 + // Korean "안녕하세요 " (6 runes each)
		250*6 // Arabic "مرحبا " (6 runes each)
	if runeCount != expectedRunes {
		t.Errorf("Rune count = %d, want %d", runeCount, expectedRunes)
	}

	// Verify UTF-8 validation
	if !te.IsValidUTF8(text) {
		t.Error("IsValidUTF8() returned false for valid large string")
	}
	if !te.IsValidUTF8Bytes(encoded) {
		t.Error("IsValidUTF8Bytes() returned false for valid large string bytes")
	}
}

// Benchmarks
func BenchmarkEncodeUTF8(b *testing.B) {
	te := &TextEncoding{}
	text := "Hello 🌍 你好"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		te.EncodeUTF8(text)
	}
}

func BenchmarkCountUTF8Bytes(b *testing.B) {
	te := &TextEncoding{}
	text := "Hello 🌍 你好"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		te.CountUTF8Bytes(text)
	}
}

func BenchmarkCountUTF8Runes(b *testing.B) {
	te := &TextEncoding{}
	text := "Hello 🌍 你好"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		te.CountUTF8Runes(text)
	}
}
