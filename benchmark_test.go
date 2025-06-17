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

func BenchmarkUTF8ByteLength(b *testing.B) {
	te := &TextEncoding{}
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " + "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " + "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		te.UTF8ByteLength(testString)
	}
}

func BenchmarkUTF8ByteLengthOptimized(b *testing.B) {
	te := &TextEncoding{}
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " + "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " + "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		te.UTF8ByteLengthOptimized(testString)
	}
}

func BenchmarkUTF8ByteLengthJS(b *testing.B) {
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " + "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " + "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utf8ByteLengthJS(testString)
	}
}

func BenchmarkUTF8ByteLengthDirect(b *testing.B) {
	testString := "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " + "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ " + "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ "

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = len(testString)
	}
}
