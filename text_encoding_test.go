package text_encoding

import (
	"encoding/base64"
	"fmt"
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

func TestInvalidUTF8Sequences(t *testing.T) {
	te := &TextEncoding{}

	// Test various invalid UTF-8 sequences
	invalidSequences := [][]byte{
		{0xFF, 0xFE},                   // Invalid start byte
		{0xC0, 0x80},                   // Overlong encoding
		{0xF0, 0x9F},                   // Incomplete sequence
		{0xED, 0xA0, 0x80},             // Surrogate pair
		{0xF4, 0x90, 0x80, 0x80},       // Out of range
		{0x80},                         // Continuation byte without start
		{0xC0, 0xAF},                   // Overlong ASCII
		{0xE0, 0x80, 0xAF},             // Overlong 2-byte sequence
		{0xF0, 0x80, 0x80, 0xAF},       // Overlong 3-byte sequence
		{0xF8, 0x80, 0x80, 0x80, 0xAF}, // 5-byte sequence (invalid)
	}

	for i, seq := range invalidSequences {
		// Test IsValidUTF8Bytes
		if te.IsValidUTF8Bytes(seq) {
			t.Errorf("IsValidUTF8Bytes should return false for invalid sequence %d", i)
		}

		// Test DecodeUTF8
		_, err := te.DecodeUTF8(seq)
		if err == nil {
			t.Errorf("DecodeUTF8 should return error for invalid sequence %d", i)
		}
	}
}

func TestConcurrentOperations(t *testing.T) {
	te := &TextEncoding{}

	// Create a test string with various characters
	testStr := "Hello 🌍 你好 café résumé 안녕하세요 مرحبا 𝄞 𒀀 👨‍👩‍👧‍👦 🏳️‍🌈"

	// Number of concurrent operations
	numGoroutines := 100

	// Channel to collect errors
	errChan := make(chan error, numGoroutines)

	// Run concurrent operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			// Encode
			encoded := te.EncodeUTF8(testStr)

			// Decode
			decoded, err := te.DecodeUTF8(encoded)
			if err != nil {
				errChan <- fmt.Errorf("decode error: %v", err)
				return
			}

			// Verify roundtrip
			if decoded != testStr {
				errChan <- fmt.Errorf("roundtrip mismatch: got %q, want %q", decoded, testStr)
				return
			}

			// Test Base64
			base64Encoded := te.EncodeUTF8ToBase64(testStr)
			base64Decoded, err := te.DecodeUTF8FromBase64(base64Encoded)
			if err != nil {
				errChan <- fmt.Errorf("base64 decode error: %v", err)
				return
			}

			if base64Decoded != testStr {
				errChan <- fmt.Errorf("base64 roundtrip mismatch: got %q, want %q", base64Decoded, testStr)
				return
			}

			errChan <- nil
		}()
	}

	// Collect errors
	for i := 0; i < numGoroutines; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("Concurrent operation failed: %v", err)
		}
	}
}

func TestStressLargeText(t *testing.T) {
	te := &TextEncoding{}

	// Create an extremely large string with various Unicode characters
	var stressText strings.Builder
	// Add a mix of characters that might stress the encoder
	for i := 0; i < 10000; i++ {
		stressText.WriteString("Hello ")
		stressText.WriteString("你好")
		stressText.WriteString("🌍")
		stressText.WriteString("café ")
		stressText.WriteString("résumé ")
		stressText.WriteString("안녕하세요 ")
		stressText.WriteString("مرحبا ")
		// Add some rare/edge case characters
		stressText.WriteString("𝄞")       // Musical symbol
		stressText.WriteString("𒀀")       // Cuneiform
		stressText.WriteString("👨‍👩‍👧‍👦") // Family emoji
		stressText.WriteString("🏳️‍🌈")    // Flag emoji
		// Add more edge cases
		stressText.WriteString("Z͑ͫ̓ͪ̂ͫ̽͏̴̙̤̞͉͚̯̞̠͍A̴̵̜̰͔ͫ͗͢L̠ͨͧͩ͘G̴̻͈͍͔̹̑͗̎̅͛́Ǫ̵̹̻̝̳͂̌̌͘!͖̬̰̙̗̿̋ͥͥ̂ͣ̐́́͜͞") // Zalgo text
		stressText.WriteString("ᚠᛇᚻ᛫ᛒᛦᚦ᛫ᚠᚱᚩᚠᚢᚱ᛫ᚠᛁᚱᚪ᛫ᚷᛖᚻᚹᛦᛚᚳᚢᛗ")                                              // Runic text
		stressText.WriteString("꧁༺༻꧂")                                                                       // Decorative characters
		stressText.WriteString("ᕕ( ᐛ )ᕗ")                                                                    // ASCII art
		stressText.WriteString("👾")                                                                          // Emoji with variation selector
		stressText.WriteString("👨‍💻")                                                                        // Emoji with ZWJ
		stressText.WriteString("🏴󠁧󠁢󠁥󠁮󠁧󠁿")                                                                    // Regional indicator
	}

	text := stressText.String()

	// Test UTF-8 encoding and decoding
	encoded := te.EncodeUTF8(text)
	decoded, err := te.DecodeUTF8(encoded)
	if err != nil {
		t.Errorf("DecodeUTF8() error: %v", err)
	}
	if decoded != text {
		t.Error("Stress test roundtrip failed")
	}

	// Test Base64 encoding and decoding
	base64Encoded := te.EncodeUTF8ToBase64(text)
	base64Decoded, err := te.DecodeUTF8FromBase64(base64Encoded)
	if err != nil {
		t.Errorf("DecodeUTF8FromBase64() error: %v", err)
	}
	if base64Decoded != text {
		t.Error("Stress test base64 roundtrip failed")
	}

	// Verify UTF-8 validation
	if !te.IsValidUTF8(text) {
		t.Error("IsValidUTF8() returned false for valid stress test string")
	}
	if !te.IsValidUTF8Bytes(encoded) {
		t.Error("IsValidUTF8Bytes() returned false for valid stress test bytes")
	}

	// Test byte and rune counting
	byteCount := te.CountUTF8Bytes(text)
	runeCount := te.CountUTF8Runes(text)

	// Verify byte count matches encoded length
	if byteCount != len(encoded) {
		t.Errorf("Byte count mismatch: CountUTF8Bytes() = %d, actual bytes length = %d", byteCount, len(encoded))
	}

	// Verify rune count is less than byte count
	if runeCount >= byteCount {
		t.Errorf("Rune count (%d) should be less than byte count (%d)", runeCount, byteCount)
	}
}

func BenchmarkTextEncoding(b *testing.B) {
	te := &TextEncoding{}

	// Create test strings of different sizes
	smallText := "Hello 🌍 你好"
	mediumText := strings.Repeat("Hello 🌍 你好 ", 100)
	largeText := strings.Repeat("Hello 🌍 你好 ", 1000)

	// Benchmark UTF-8 encoding
	b.Run("EncodeUTF8-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.EncodeUTF8(smallText)
		}
	})

	b.Run("EncodeUTF8-Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.EncodeUTF8(mediumText)
		}
	})

	b.Run("EncodeUTF8-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.EncodeUTF8(largeText)
		}
	})

	// Benchmark UTF-8 decoding
	smallEncoded := te.EncodeUTF8(smallText)
	mediumEncoded := te.EncodeUTF8(mediumText)
	largeEncoded := te.EncodeUTF8(largeText)

	b.Run("DecodeUTF8-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.DecodeUTF8(smallEncoded)
		}
	})

	b.Run("DecodeUTF8-Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.DecodeUTF8(mediumEncoded)
		}
	})

	b.Run("DecodeUTF8-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.DecodeUTF8(largeEncoded)
		}
	})

	// Benchmark Base64 encoding
	b.Run("EncodeUTF8ToBase64-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.EncodeUTF8ToBase64(smallText)
		}
	})

	b.Run("EncodeUTF8ToBase64-Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.EncodeUTF8ToBase64(mediumText)
		}
	})

	b.Run("EncodeUTF8ToBase64-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.EncodeUTF8ToBase64(largeText)
		}
	})

	// Benchmark Base64 decoding
	smallBase64 := te.EncodeUTF8ToBase64(smallText)
	mediumBase64 := te.EncodeUTF8ToBase64(mediumText)
	largeBase64 := te.EncodeUTF8ToBase64(largeText)

	b.Run("DecodeUTF8FromBase64-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.DecodeUTF8FromBase64(smallBase64)
		}
	})

	b.Run("DecodeUTF8FromBase64-Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.DecodeUTF8FromBase64(mediumBase64)
		}
	})

	b.Run("DecodeUTF8FromBase64-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.DecodeUTF8FromBase64(largeBase64)
		}
	})

	// Benchmark counting functions
	b.Run("CountUTF8Bytes-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.CountUTF8Bytes(smallText)
		}
	})

	b.Run("CountUTF8Runes-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.CountUTF8Runes(smallText)
		}
	})

	b.Run("IsValidUTF8-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.IsValidUTF8(smallText)
		}
	})

	b.Run("IsValidUTF8Bytes-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			te.IsValidUTF8Bytes(smallEncoded)
		}
	})
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
