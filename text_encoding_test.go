package text_encoding

import (
	"testing"
)

func TestGetEncoding(t *testing.T) {
	tests := []struct {
		label    string
		expected bool
	}{
		{"utf-8", true},
		{"UTF-8", true},
		{"utf8", true},
		{"iso-8859-1", true},
		{"latin1", true},
		{"shift-jis", true},
		{"sjis", true},
		{"invalid-encoding", false},
		{"", false},
	}

	for _, test := range tests {
		_, err := getEncoding(test.label)
		if test.expected && err != nil {
			t.Errorf("getEncoding(%q) should succeed, got error: %v", test.label, err)
		}
		if !test.expected && err == nil {
			t.Errorf("getEncoding(%q) should fail, but succeeded", test.label)
		}
	}
}

func TestTextEncoder(t *testing.T) {
	te := &TextEncoding{}

	// Test UTF-8 encoder
	encoder, err := te.NewTextEncoder("utf-8")
	if err != nil {
		t.Fatalf("Failed to create UTF-8 encoder: %v", err)
	}

	// Test basic encoding
	text := "Hello, 世界!"
	encoded, err := encoder.Encode(text)
	if err != nil {
		t.Fatalf("Failed to encode text: %v", err)
	}

	if len(encoded) == 0 {
		t.Error("Encoded data should not be empty")
	}

	// Test empty string
	emptyEncoded, err := encoder.Encode("")
	if err != nil {
		t.Fatalf("Failed to encode empty string: %v", err)
	}

	if len(emptyEncoded) != 0 {
		t.Error("Empty string should encode to empty bytes")
	}
}

func TestTextDecoder(t *testing.T) {
	te := &TextEncoding{}

	// Test UTF-8 decoder
	decoder, err := te.NewTextDecoder("utf-8")
	if err != nil {
		t.Fatalf("Failed to create UTF-8 decoder: %v", err)
	}

	// Test basic decoding
	originalText := "Hello, 世界!"
	encoded := []byte(originalText)
	decoded, err := decoder.Decode(encoded)
	if err != nil {
		t.Fatalf("Failed to decode data: %v", err)
	}

	if decoded != originalText {
		t.Errorf("Decoded text doesn't match original. Expected: %q, Got: %q", originalText, decoded)
	}

	// Test empty bytes
	emptyDecoded, err := decoder.Decode([]byte{})
	if err != nil {
		t.Fatalf("Failed to decode empty bytes: %v", err)
	}

	if emptyDecoded != "" {
		t.Error("Empty bytes should decode to empty string")
	}
}

func TestTextEncodingRoundTrip(t *testing.T) {
	te := &TextEncoding{}

	// Test UTF-8 round trip
	encoder, err := te.NewTextEncoder("utf-8")
	if err != nil {
		t.Fatalf("Failed to create encoder: %v", err)
	}

	decoder, err := te.NewTextDecoder("utf-8")
	if err != nil {
		t.Fatalf("Failed to create decoder: %v", err)
	}

	testCases := []string{
		"Hello, World!",
		"Hello, 世界!",
		"Special chars: áéíóú ñ ç ß € ¥ £",
		"",
		"1234567890",
		"!@#$%^&*()_+-=[]{}|;':\",./<>?",
	}

	for _, testCase := range testCases {
		encoded, err := encoder.Encode(testCase)
		if err != nil {
			t.Fatalf("Failed to encode %q: %v", testCase, err)
		}

		decoded, err := decoder.Decode(encoded)
		if err != nil {
			t.Fatalf("Failed to decode %q: %v", testCase, err)
		}

		if decoded != testCase {
			t.Errorf("Round trip failed for %q. Expected: %q, Got: %q", testCase, testCase, decoded)
		}
	}
}

func TestIsValidEncoding(t *testing.T) {
	te := &TextEncoding{}

	tests := []struct {
		label    string
		expected bool
	}{
		{"utf-8", true},
		{"iso-8859-1", true},
		{"shift-jis", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := te.IsValidEncoding(test.label)
		if result != test.expected {
			t.Errorf("IsValidEncoding(%q) = %v, expected %v", test.label, result, test.expected)
		}
	}
}

func TestGetSupportedEncodings(t *testing.T) {
	te := &TextEncoding{}
	supported := te.GetSupportedEncodings()

	if len(supported) == 0 {
		t.Error("GetSupportedEncodings should return non-empty list")
	}

	// Check that some expected encodings are present
	expectedEncodings := []string{"utf-8", "iso-8859-1", "shift-jis"}
	for _, expected := range expectedEncodings {
		found := false
		for _, supported := range supported {
			if supported == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected encoding %q not found in supported encodings", expected)
		}
	}
}
