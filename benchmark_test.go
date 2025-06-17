package text_encoding

import (
	"testing"
	"unicode/utf8"
)

// Simulate JavaScript implementation for benchmarking
func utf8ByteLengthJS(str string) int {
	bytes := 0
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		if r == utf8.RuneError {
			// Handle invalid UTF-8
			bytes += 1
			i++
		} else {
			// Calculate bytes based on Unicode code point
			if r <= 0x7f {
				bytes += 1
			} else if r <= 0x7ff {
				bytes += 2
			} else if r <= 0xffff {
				bytes += 3
			} else {
				bytes += 4
			}
			i += size
		}
	}
	return bytes
}

// Manual UTF-8 byte length calculation (educational/alternative implementation)
func utf8ByteLengthManual(str string) int {
	bytes := 0
	for _, r := range str {
		switch {
		case r <= 0x7f:
			bytes += 1
		case r <= 0x7ff:
			bytes += 2
		case r <= 0xffff:
			bytes += 3
		default:
			bytes += 4
		}
	}
	return bytes
}

// Benchmark our Utils.UTF8ByteLength method
func BenchmarkUTF8ByteLength(b *testing.B) {
	utils := &Utils{}
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.UTF8ByteLength(testString)
	}
}

// Benchmark JavaScript-like implementation
func BenchmarkUTF8ByteLengthJS(b *testing.B) {
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utf8ByteLengthJS(testString)
	}
}

// Benchmark manual calculation (rune iteration)
func BenchmarkUTF8ByteLengthManual(b *testing.B) {
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utf8ByteLengthManual(testString)
	}
}

// Benchmark direct len() call (Go's optimized UTF-8 handling)
func BenchmarkUTF8ByteLengthDirect(b *testing.B) {
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = len(testString)
	}
}

// Benchmark using utf8.RuneCountInString (counts runes, not bytes)
func BenchmarkUTF8RuneCount(b *testing.B) {
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " +
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utf8.RuneCountInString(testString)
	}
}

// Benchmark different string sizes
func BenchmarkUTF8ByteLength_Small(b *testing.B) {
	utils := &Utils{}
	testString := "Hello"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.UTF8ByteLength(testString)
	}
}

func BenchmarkUTF8ByteLength_Medium(b *testing.B) {
	utils := &Utils{}
	testString := "Hello, 世界! 🌍"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.UTF8ByteLength(testString)
	}
}

func BenchmarkUTF8ByteLength_Large(b *testing.B) {
	utils := &Utils{}
	baseString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "
	testString := ""
	for i := 0; i < 100; i++ {
		testString += baseString
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.UTF8ByteLength(testString)
	}
}

// Benchmark ASCII-only strings
func BenchmarkUTF8ByteLength_ASCII(b *testing.B) {
	utils := &Utils{}
	testString := "Hello World! This is a test string with only ASCII characters."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.UTF8ByteLength(testString)
	}
}

// Benchmark Unicode-heavy strings
func BenchmarkUTF8ByteLength_Unicode(b *testing.B) {
	utils := &Utils{}
	testString := "世界🌍世界🌍世界🌍世界🌍世界🌍世界🌍世界🌍"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.UTF8ByteLength(testString)
	}
}

// Benchmark empty string
func BenchmarkUTF8ByteLength_Empty(b *testing.B) {
	utils := &Utils{}
	testString := ""

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.UTF8ByteLength(testString)
	}
}

// Test to verify all implementations give the same result
func TestUTF8ByteLengthConsistency(t *testing.T) {
	utils := &Utils{}

	testCases := []string{
		"",
		"Hello",
		"Hello, 世界!",
		"🌍",
		"áéíóú ñ ç ß € ¥ £",
		"Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £",
	}

	for _, testString := range testCases {
		t.Run(testString, func(t *testing.T) {
			expected := len(testString) // Go's built-in UTF-8 byte length

			result1 := utils.UTF8ByteLength(testString)
			result2 := utf8ByteLengthJS(testString)
			result3 := utf8ByteLengthManual(testString)

			if result1 != expected {
				t.Errorf("Utils.UTF8ByteLength(%q) = %d, want %d", testString, result1, expected)
			}
			if result2 != expected {
				t.Errorf("utf8ByteLengthJS(%q) = %d, want %d", testString, result2, expected)
			}
			if result3 != expected {
				t.Errorf("utf8ByteLengthManual(%q) = %d, want %d", testString, result3, expected)
			}
		})
	}
}
