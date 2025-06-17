package text_encoding

import (
	"encoding/base64"
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
			name:        "nil input",
			input:       nil,
			expected:    "",
			expectError: true,
		},
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
			expected: 12, // "Hello " (6 bytes) + 🌍 (4 bytes) = 10 bytes
		},
		{
			name:     "chinese characters",
			input:    "你好",
			expected: 6, // 3 bytes per character
		},
		{
			name:     "mixed content",
			input:    "café 🚀",
			expected: 10, // c(1) + a(1) + f(1) + é(2) + space(1) + rocket(4) = 10
		},
		{
			name:     "only emojis",
			input:    "🚀🌍💻",
			expected: 12, // 4 bytes each
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.CountUTF8Bytes(tt.input)
			if result != tt.expected {
				t.Errorf("CountUTF8Bytes() = %d, want %d", result, tt.expected)
				t.Errorf("Actual byte length: %d", len([]byte(tt.input)))
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
			expected: 7, // 6 chars + 1 emoji
		},
		{
			name:     "chinese characters",
			input:    "你好",
			expected: 2,
		},
		{
			name:     "mixed content",
			input:    "café 🚀",
			expected: 6, // c + a + f + é + space + rocket
		},
		{
			name:     "only emojis",
			input:    "🚀🌍💻",
			expected: 3, // 3 emoji characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := te.CountUTF8Runes(tt.input)
			if result != tt.expected {
				t.Errorf("CountUTF8Runes() = %d, want %d", result, tt.expected)
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
			input:    "Hello 🌍 你好",
			expected: true,
		},
		{
			name:     "valid string with special chars",
			input:    "café naïve résumé",
			expected: true,
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
			name:     "nil input",
			input:    nil,
			expected: true, // utf8.Valid returns true for nil
		},
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: true,
		},
		{
			name:     "valid ascii bytes",
			input:    []byte("hello"),
			expected: true,
		},
		{
			name:     "valid unicode bytes",
			input:    []byte("Hello 🌍"),
			expected: true,
		},
		{
			name:     "invalid utf-8 sequence",
			input:    []byte{0xFF, 0xFE},
			expected: false,
		},
		{
			name:     "incomplete utf-8 sequence",
			input:    []byte{0xF0, 0x9F},
			expected: false,
		},
		{
			name:     "overlong encoding",
			input:    []byte{0xC0, 0x80}, // overlong encoding of null byte
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

// Benchmark tests
func BenchmarkEncodeUTF8(b *testing.B) {
	te := &TextEncoding{}
	text := "Hello 🌍 世界 café naïve résumé"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		te.EncodeUTF8(text)
	}
}

func BenchmarkCountUTF8Bytes(b *testing.B) {
	te := &TextEncoding{}
	text := "Hello 🌍 世界 café naïve résumé"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		te.CountUTF8Bytes(text)
	}
}

func BenchmarkCountUTF8Runes(b *testing.B) {
	te := &TextEncoding{}
	text := "Hello 🌍 世界 café naïve résumé"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		te.CountUTF8Runes(text)
	}
}
