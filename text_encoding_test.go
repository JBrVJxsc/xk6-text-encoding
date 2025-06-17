package text_encoding

import (
	"bytes"
	"testing"

	"golang.org/x/text/encoding/unicode"
)

func TestGetEncoding(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		wantErr  bool
		expected string
	}{
		{"UTF-8", "utf-8", false, "utf-8"},
		{"UTF-8 uppercase", "UTF-8", false, "utf-8"},
		{"UTF-8 with spaces", " utf-8 ", false, "utf-8"},
		{"UTF-16", "utf-16", false, "utf-16"},
		{"ISO-8859-1", "iso-8859-1", false, "iso-8859-1"},
		{"Latin1 alias", "latin1", false, "iso-8859-1"},
		{"Windows-1252", "windows-1252", false, "windows-1252"},
		{"Shift-JIS", "shift-jis", false, "shift-jis"},
		{"Shift-JIS variant", "sjis", false, "shift-jis"},
		{"GBK", "gbk", false, "gbk"},
		{"Invalid encoding", "invalid-encoding", true, ""},
		{"Empty string", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoding, err := getEncoding(tt.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEncoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && encoding == nil {
				t.Errorf("getEncoding() returned nil encoding for valid label")
			}
		})
	}
}

func TestTextEncoder_Encode(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		text     string
		wantErr  bool
		expected []byte
	}{
		{"UTF-8 simple", "utf-8", "Hello", false, []byte("Hello")},
		{"UTF-8 unicode", "utf-8", "Hello, 世界!", false, []byte("Hello, 世界!")},
		{"UTF-8 emoji", "utf-8", "🌍", false, []byte("🌍")},
		{"UTF-8 empty", "utf-8", "", false, []byte{}},
		{"ISO-8859-1 simple", "iso-8859-1", "Hello", false, nil}, // We'll check length
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := &TextEncoder{label: tt.label}
			enc, err := getEncoding(tt.label)
			if err != nil {
				t.Fatalf("Failed to get encoding: %v", err)
			}
			encoder.encoding = enc

			result, err := encoder.Encode(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("TextEncoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.expected != nil && !bytes.Equal(result, tt.expected) {
					t.Errorf("TextEncoder.Encode() = %v, want %v", result, tt.expected)
				}
				if len(result) == 0 && tt.text != "" {
					t.Errorf("TextEncoder.Encode() returned empty result for non-empty input")
				}
			}
		})
	}
}

func TestTextEncoder_EncodeString(t *testing.T) {
	encoder := &TextEncoder{
		encoding: unicode.UTF8,
		label:    "utf-8",
	}

	tests := []struct {
		name string
		text string
	}{
		{"Simple text", "Hello World"},
		{"Unicode text", "Hello, 世界!"},
		{"Empty string", ""},
		{"Special chars", "áéíóú ñ ç ß € ¥ £"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoder.EncodeString(tt.text)
			if err != nil {
				t.Errorf("TextEncoder.EncodeString() error = %v", err)
				return
			}

			// For UTF-8, the string should be the same as the original
			if encoder.label == "utf-8" && result != tt.text {
				t.Errorf("TextEncoder.EncodeString() = %v, want %v", result, tt.text)
			}
		})
	}
}

func TestTextEncoder_GetEncoding(t *testing.T) {
	encoder := &TextEncoder{label: "utf-8"}
	if got := encoder.GetEncoding(); got != "utf-8" {
		t.Errorf("TextEncoder.GetEncoding() = %v, want %v", got, "utf-8")
	}
}

func TestTextDecoder_Decode(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		data     []byte
		expected string
		wantErr  bool
	}{
		{"UTF-8 simple", "utf-8", []byte("Hello"), "Hello", false},
		{"UTF-8 unicode", "utf-8", []byte("Hello, 世界!"), "Hello, 世界!", false},
		{"UTF-8 empty", "utf-8", []byte{}, "", false},
		{"UTF-8 emoji", "utf-8", []byte("🌍"), "🌍", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := &TextDecoder{label: tt.label}
			enc, err := getEncoding(tt.label)
			if err != nil {
				t.Fatalf("Failed to get encoding: %v", err)
			}
			decoder.encoding = enc

			result, err := decoder.Decode(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TextDecoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected {
				t.Errorf("TextDecoder.Decode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTextDecoder_GetEncoding(t *testing.T) {
	decoder := &TextDecoder{label: "iso-8859-1"}
	if got := decoder.GetEncoding(); got != "iso-8859-1" {
		t.Errorf("TextDecoder.GetEncoding() = %v, want %v", got, "iso-8859-1")
	}
}

func TestUtils_UTF8ByteLength(t *testing.T) {
	utils := &Utils{}

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"ASCII only", "Hello", 5},
		{"ASCII with space", "Hello World", 11},
		{"Unicode chars", "Hello, 世界!", 13}, // 7 ASCII + 6 bytes for 2 Chinese chars
		{"Emoji", "🌍", 4},                   // Emoji is 4 bytes in UTF-8
		{"Mixed", "Hello 🌍 世界", 16},         // 6 ASCII + 4 emoji + 1 space + 6 Chinese + 1 space
		{"Special chars", "áéíóú", 10},      // Each accented char is 2 bytes in UTF-8
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.UTF8ByteLength(tt.input)
			if result != tt.expected {
				t.Errorf("Utils.UTF8ByteLength(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUtils_IsValidEncoding(t *testing.T) {
	utils := &Utils{}

	tests := []struct {
		name     string
		encoding string
		expected bool
	}{
		{"UTF-8", "utf-8", true},
		{"UTF-8 uppercase", "UTF-8", true},
		{"UTF-16", "utf-16", true},
		{"ISO-8859-1", "iso-8859-1", true},
		{"Latin1 alias", "latin1", true},
		{"Windows-1252", "windows-1252", true},
		{"Shift-JIS", "shift-jis", true},
		{"GBK", "gbk", true},
		{"Invalid", "invalid-encoding", false},
		{"Empty", "", false},
		{"Nonsense", "definitely-not-real", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IsValidEncoding(tt.encoding)
			if result != tt.expected {
				t.Errorf("Utils.IsValidEncoding(%q) = %v, want %v", tt.encoding, result, tt.expected)
			}
		})
	}
}

func TestUtils_GetSupportedEncodings(t *testing.T) {
	utils := &Utils{}

	encodings := utils.GetSupportedEncodings()

	if len(encodings) == 0 {
		t.Error("GetSupportedEncodings() returned empty slice")
	}

	// Check that some expected encodings are in the list
	expectedEncodings := []string{"utf-8", "utf-16", "iso-8859-1", "windows-1252", "shift-jis", "gbk"}

	for _, expected := range expectedEncodings {
		found := false
		for _, encoding := range encodings {
			if encoding == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected encoding %q not found in supported encodings", expected)
		}
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	tests := []struct {
		name     string
		encoding string
		text     string
	}{
		{"UTF-8", "utf-8", "Hello, 世界! 🌍"},
		{"ISO-8859-1", "iso-8859-1", "Hello World"},
		{"Windows-1252", "windows-1252", "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := getEncoding(tt.encoding)
			if err != nil {
				t.Fatalf("Failed to get encoding: %v", err)
			}

			encoder := &TextEncoder{encoding: enc, label: tt.encoding}
			decoder := &TextDecoder{encoding: enc, label: tt.encoding}

			// Encode
			encoded, err := encoder.Encode(tt.text)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}

			// Decode
			decoded, err := decoder.Decode(encoded)
			if err != nil {
				t.Fatalf("Failed to decode: %v", err)
			}

			// For encodings that can't represent all Unicode chars, we might not get exact match
			// But for UTF-8, we should get exact match
			if tt.encoding == "utf-8" && decoded != tt.text {
				t.Errorf("Roundtrip failed: got %q, want %q", decoded, tt.text)
			}
		})
	}
}

func TestBufferPool(t *testing.T) {
	// Test that buffer pool works correctly
	buf1 := getBuffer()
	buf1.WriteString("test")

	if buf1.Len() != 4 {
		t.Errorf("Buffer should have length 4, got %d", buf1.Len())
	}

	putBuffer(buf1)

	// After putting back, buffer should be reset
	buf2 := getBuffer()
	if buf2.Len() != 0 {
		t.Errorf("Buffer should be reset, got length %d", buf2.Len())
	}

	putBuffer(buf2)
}
