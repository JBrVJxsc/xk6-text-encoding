package text_encoding

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestEncodeUTF8(t *testing.T) {
	te := &TextEncoding{}

	// Empty string
	result, err := te.EncodeUTF8("")
	if err != nil {
		t.Errorf("EncodeUTF8() error: %v", err)
	}
	if len(result) != 0 {
		t.Error("Empty string should produce empty bytes")
	}

	// ASCII text
	result, err = te.EncodeUTF8("hello")
	if err != nil {
		t.Errorf("EncodeUTF8() error: %v", err)
	}
	if len(result) != 5 {
		t.Error("ASCII 'hello' should be 5 bytes")
	}
	expected := []byte{104, 101, 108, 108, 111}
	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Byte %d: got %d, want %d", i, b, expected[i])
		}
	}

	// Unicode text with emoji
	result, err = te.EncodeUTF8("Hello 🌍")
	if err != nil {
		t.Errorf("EncodeUTF8() error: %v", err)
	}
	if len(result) != 10 {
		t.Error("Unicode with emoji should be 10 bytes")
	}

	// Chinese characters
	result, err = te.EncodeUTF8("你好")
	if err != nil {
		t.Errorf("EncodeUTF8() error: %v", err)
	}
	if len(result) != 6 {
		t.Error("Chinese characters should be 6 bytes (3 each)")
	}
}

func TestEncodeUTF8ToBase64(t *testing.T) {
	te := &TextEncoding{}

	// Empty string
	result, err := te.EncodeUTF8ToBase64("")
	if err != nil {
		t.Errorf("EncodeUTF8ToBase64() error: %v", err)
	}
	if result != "" {
		t.Error("Empty string should produce empty base64")
	}

	// Simple text
	result, err = te.EncodeUTF8ToBase64("hello")
	if err != nil {
		t.Errorf("EncodeUTF8ToBase64() error: %v", err)
	}
	if result == "" {
		t.Error("Base64 result should not be empty")
	}

	// Verify round-trip
	decoded, err := te.DecodeUTF8FromBase64(result)
	if err != nil {
		t.Errorf("DecodeUTF8FromBase64() error: %v", err)
	}
	if decoded != "hello" {
		t.Error("Round-trip base64 should work")
	}

	// Unicode text
	result, err = te.EncodeUTF8ToBase64("Hello 🌍")
	if err != nil {
		t.Errorf("EncodeUTF8ToBase64() error: %v", err)
	}
	decoded, err = te.DecodeUTF8FromBase64(result)
	if err != nil {
		t.Errorf("DecodeUTF8FromBase64() error: %v", err)
	}
	if decoded != "Hello 🌍" {
		t.Error("Unicode base64 round-trip should work")
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

	count, err := te.CountUTF8Bytes("")
	if err != nil {
		t.Errorf("CountUTF8Bytes() error: %v", err)
	}
	if count != 0 {
		t.Error("Empty string should have 0 bytes")
	}

	count, err = te.CountUTF8Bytes("hello")
	if err != nil {
		t.Errorf("CountUTF8Bytes() error: %v", err)
	}
	if count != 5 {
		t.Error("ASCII should have 5 bytes")
	}

	count, err = te.CountUTF8Bytes("Hello 🌍")
	if err != nil {
		t.Errorf("CountUTF8Bytes() error: %v", err)
	}
	if count != 10 {
		t.Error("Unicode with emoji should have 10 bytes")
	}

	count, err = te.CountUTF8Bytes("你好")
	if err != nil {
		t.Errorf("CountUTF8Bytes() error: %v", err)
	}
	if count != 6 {
		t.Error("Chinese characters should have 6 bytes")
	}
}

func TestCountUTF8Runes(t *testing.T) {
	te := &TextEncoding{}

	count, err := te.CountUTF8Runes("")
	if err != nil {
		t.Errorf("CountUTF8Runes() error: %v", err)
	}
	if count != 0 {
		t.Error("Empty string should have 0 runes")
	}

	count, err = te.CountUTF8Runes("hello")
	if err != nil {
		t.Errorf("CountUTF8Runes() error: %v", err)
	}
	if count != 5 {
		t.Error("ASCII should have 5 runes")
	}

	count, err = te.CountUTF8Runes("Hello 🌍")
	if err != nil {
		t.Errorf("CountUTF8Runes() error: %v", err)
	}
	if count != 7 {
		t.Error("Unicode with emoji should have 7 runes")
	}

	count, err = te.CountUTF8Runes("你好")
	if err != nil {
		t.Errorf("CountUTF8Runes() error: %v", err)
	}
	if count != 2 {
		t.Error("Chinese characters should have 2 runes")
	}
}

func TestIsValidUTF8(t *testing.T) {
	te := &TextEncoding{}

	valid, err := te.IsValidUTF8("")
	if err != nil {
		t.Errorf("IsValidUTF8() error: %v", err)
	}
	if !valid {
		t.Error("Empty string should be valid")
	}

	valid, err = te.IsValidUTF8("hello")
	if err != nil {
		t.Errorf("IsValidUTF8() error: %v", err)
	}
	if !valid {
		t.Error("ASCII should be valid")
	}

	valid, err = te.IsValidUTF8("Hello 🌍 你好")
	if err != nil {
		t.Errorf("IsValidUTF8() error: %v", err)
	}
	if !valid {
		t.Error("Unicode should be valid")
	}

	valid, err = te.IsValidUTF8("café naïve résumé")
	if err != nil {
		t.Errorf("IsValidUTF8() error: %v", err)
	}
	if !valid {
		t.Error("Special chars should be valid")
	}

	valid, err = te.IsValidUTF8("🚀🌍💻中文한국어العربية")
	if err != nil {
		t.Errorf("IsValidUTF8() error: %v", err)
	}
	if !valid {
		t.Error("Complex Unicode should be valid")
	}
}

func TestIsValidUTF8Bytes(t *testing.T) {
	te := &TextEncoding{}

	valid, err := te.IsValidUTF8Bytes([]byte{})
	if err != nil {
		t.Errorf("IsValidUTF8Bytes() error: %v", err)
	}
	if !valid {
		t.Error("Empty bytes should be valid")
	}

	bytes := []byte{104, 101, 108, 108, 111} // 'hello'
	valid, err = te.IsValidUTF8Bytes(bytes)
	if err != nil {
		t.Errorf("IsValidUTF8Bytes() error: %v", err)
	}
	if !valid {
		t.Error("ASCII bytes should be valid")
	}

	bytes = []byte("Hello 🌍")
	valid, err = te.IsValidUTF8Bytes(bytes)
	if err != nil {
		t.Errorf("IsValidUTF8Bytes() error: %v", err)
	}
	if !valid {
		t.Error("Unicode bytes should be valid")
	}

	invalidBytes := []byte{0xFF, 0xFE}
	valid, err = te.IsValidUTF8Bytes(invalidBytes)
	if err != nil {
		t.Errorf("IsValidUTF8Bytes() error: %v", err)
	}
	if valid {
		t.Error("Invalid UTF-8 bytes should not be valid")
	}

	incompleteBytes := []byte{0xF0, 0x9F}
	valid, err = te.IsValidUTF8Bytes(incompleteBytes)
	if err != nil {
		t.Errorf("IsValidUTF8Bytes() error: %v", err)
	}
	if valid {
		t.Error("Incomplete UTF-8 should not be valid")
	}

	overlongBytes := []byte{0xC0, 0x80}
	valid, err = te.IsValidUTF8Bytes(overlongBytes)
	if err != nil {
		t.Errorf("IsValidUTF8Bytes() error: %v", err)
	}
	if valid {
		t.Error("Overlong encoding should not be valid")
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
		valid, err := te.IsValidUTF8Bytes(seq)
		if err != nil {
			t.Errorf("IsValidUTF8Bytes() error for sequence %d: %v", i, err)
		}
		if valid {
			t.Errorf("IsValidUTF8Bytes should return false for invalid sequence %d", i)
		}

		// Test DecodeUTF8
		_, err = te.DecodeUTF8(seq)
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
			encoded, err := te.EncodeUTF8(testStr)
			if err != nil {
				errChan <- fmt.Errorf("encode error: %v", err)
				return
			}

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
			base64Encoded, err := te.EncodeUTF8ToBase64(testStr)
			if err != nil {
				errChan <- fmt.Errorf("base64 encode error: %v", err)
				return
			}

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
	encoded, err := te.EncodeUTF8(text)
	if err != nil {
		t.Errorf("EncodeUTF8() error: %v", err)
	}
	decoded, err := te.DecodeUTF8(encoded)
	if err != nil {
		t.Errorf("DecodeUTF8() error: %v", err)
	}
	if decoded != text {
		t.Error("Stress test roundtrip failed")
	}

	// Test Base64 encoding and decoding
	base64Encoded, err := te.EncodeUTF8ToBase64(text)
	if err != nil {
		t.Errorf("EncodeUTF8ToBase64() error: %v", err)
	}
	base64Decoded, err := te.DecodeUTF8FromBase64(base64Encoded)
	if err != nil {
		t.Errorf("DecodeUTF8FromBase64() error: %v", err)
	}
	if base64Decoded != text {
		t.Error("Stress test base64 roundtrip failed")
	}

	// Verify UTF-8 validation
	valid, err := te.IsValidUTF8(text)
	if err != nil {
		t.Errorf("IsValidUTF8() error: %v", err)
	}
	if !valid {
		t.Error("IsValidUTF8() returned false for valid stress test string")
	}

	valid, err = te.IsValidUTF8Bytes(encoded)
	if err != nil {
		t.Errorf("IsValidUTF8Bytes() error: %v", err)
	}
	if !valid {
		t.Error("IsValidUTF8Bytes() returned false for valid stress test bytes")
	}

	// Test byte and rune counting
	byteCount, err := te.CountUTF8Bytes(text)
	if err != nil {
		t.Errorf("CountUTF8Bytes() error: %v", err)
	}
	runeCount, err := te.CountUTF8Runes(text)
	if err != nil {
		t.Errorf("CountUTF8Runes() error: %v", err)
	}

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

	// Create test data of different sizes
	smallText := "Hello 🌍"
	mediumText := strings.Repeat("Hello 🌍 你好 ", 100)
	largeText := strings.Repeat("Hello 🌍 你好 ", 1000)

	// Benchmark UTF-8 encoding
	b.Run("EncodeUTF8-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.EncodeUTF8(smallText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("EncodeUTF8-Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.EncodeUTF8(mediumText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("EncodeUTF8-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.EncodeUTF8(largeText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark UTF-8 decoding
	smallEncoded, err := te.EncodeUTF8(smallText)
	if err != nil {
		b.Fatal(err)
	}
	mediumEncoded, err := te.EncodeUTF8(mediumText)
	if err != nil {
		b.Fatal(err)
	}
	largeEncoded, err := te.EncodeUTF8(largeText)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("DecodeUTF8-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.DecodeUTF8(smallEncoded)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeUTF8-Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.DecodeUTF8(mediumEncoded)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeUTF8-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.DecodeUTF8(largeEncoded)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark Base64 encoding
	b.Run("EncodeUTF8ToBase64-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.EncodeUTF8ToBase64(smallText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("EncodeUTF8ToBase64-Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.EncodeUTF8ToBase64(mediumText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("EncodeUTF8ToBase64-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.EncodeUTF8ToBase64(largeText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark Base64 decoding
	smallBase64, err := te.EncodeUTF8ToBase64(smallText)
	if err != nil {
		b.Fatal(err)
	}
	mediumBase64, err := te.EncodeUTF8ToBase64(mediumText)
	if err != nil {
		b.Fatal(err)
	}
	largeBase64, err := te.EncodeUTF8ToBase64(largeText)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("DecodeUTF8FromBase64-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.DecodeUTF8FromBase64(smallBase64)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeUTF8FromBase64-Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.DecodeUTF8FromBase64(mediumBase64)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DecodeUTF8FromBase64-Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.DecodeUTF8FromBase64(largeBase64)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark counting functions
	b.Run("CountUTF8Bytes-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.CountUTF8Bytes(smallText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("CountUTF8Runes-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.CountUTF8Runes(smallText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("IsValidUTF8-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.IsValidUTF8(smallText)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("IsValidUTF8Bytes-Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := te.IsValidUTF8Bytes(smallEncoded)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEncodeUTF8(b *testing.B) {
	te := &TextEncoding{}
	testCases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"ascii", "hello"},
		{"unicode", "Hello 🌍"},
		{"chinese", "你好"},
		{"mixed", "Hello 🌍 你好"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := te.EncodeUTF8(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecodeUTF8(b *testing.B) {
	te := &TextEncoding{}
	testCases := []struct {
		name  string
		input []byte
	}{
		{"empty", []byte{}},
		{"ascii", []byte("hello")},
		{"unicode", []byte("Hello 🌍")},
		{"chinese", []byte("你好")},
		{"mixed", []byte("Hello 🌍 你好")},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := te.DecodeUTF8(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEncodeUTF8ToBase64(b *testing.B) {
	te := &TextEncoding{}
	testCases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"ascii", "hello"},
		{"unicode", "Hello 🌍"},
		{"chinese", "你好"},
		{"mixed", "Hello 🌍 你好"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := te.EncodeUTF8ToBase64(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecodeUTF8FromBase64(b *testing.B) {
	te := &TextEncoding{}
	testCases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"ascii", base64.StdEncoding.EncodeToString([]byte("hello"))},
		{"unicode", base64.StdEncoding.EncodeToString([]byte("Hello 🌍"))},
		{"chinese", base64.StdEncoding.EncodeToString([]byte("你好"))},
		{"mixed", base64.StdEncoding.EncodeToString([]byte("Hello 🌍 你好"))},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := te.DecodeUTF8FromBase64(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCountUTF8Bytes(b *testing.B) {
	te := &TextEncoding{}
	testCases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"ascii", "hello"},
		{"unicode", "Hello 🌍"},
		{"chinese", "你好"},
		{"mixed", "Hello 🌍 你好"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := te.CountUTF8Bytes(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCountUTF8Runes(b *testing.B) {
	te := &TextEncoding{}
	testCases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"ascii", "hello"},
		{"unicode", "Hello 🌍"},
		{"chinese", "你好"},
		{"mixed", "Hello 🌍 你好"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := te.CountUTF8Runes(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkIsValidUTF8(b *testing.B) {
	te := &TextEncoding{}
	testCases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"ascii", "hello"},
		{"unicode", "Hello 🌍"},
		{"chinese", "你好"},
		{"mixed", "Hello 🌍 你好"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := te.IsValidUTF8(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkIsValidUTF8Bytes(b *testing.B) {
	te := &TextEncoding{}
	testCases := []struct {
		name  string
		input []byte
	}{
		{"empty", []byte{}},
		{"ascii", []byte("hello")},
		{"unicode", []byte("Hello 🌍")},
		{"chinese", []byte("你好")},
		{"mixed", []byte("Hello 🌍 你好")},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := te.IsValidUTF8Bytes(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
