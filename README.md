# k6 Text Encoding Extension

A k6 extension for text encoding and decoding operations, providing high-performance UTF-8 and Base64 handling.

## Installation

```bash
go install go.k6.io/xk6/cmd/xk6@latest
xk6 build --with github.com/JBrVJxsc/xk6-text-encoding@latest
```

## Usage

Import the module in your k6 script:

```javascript
import encoding from 'k6/x/text-encoding';
```

### Basic Operations

#### UTF-8 Encoding and Decoding

```javascript
// Encode string to UTF-8 bytes
const bytes = encoding.encodeUTF8('Hello 🌍');
console.log(bytes.length); // 10 bytes

// Decode UTF-8 bytes to string
const text = encoding.decodeUTF8(bytes);
console.log(text); // "Hello 🌍"

// Handle empty strings
const emptyBytes = encoding.encodeUTF8('');
console.log(emptyBytes.length); // 0
```

#### Base64 Encoding and Decoding

```javascript
// Encode string to Base64
const base64 = encoding.encodeUTF8ToBase64('Hello 🌍');
console.log(base64); // Base64 encoded string

// Decode Base64 to string
const decoded = encoding.decodeUTF8FromBase64(base64);
console.log(decoded); // "Hello 🌍"

// Round-trip verification
const original = 'café naïve résumé 中文 🚀🌍💻';
const encoded = encoding.encodeUTF8ToBase64(original);
const roundtrip = encoding.decodeUTF8FromBase64(encoded);
console.log(roundtrip === original); // true
```

#### UTF-8 Validation

```javascript
// Check if string is valid UTF-8
const isValid = encoding.isValidUTF8('Hello 🌍 你好');
console.log(isValid); // true

// Check if bytes are valid UTF-8
const bytes = new Uint8Array([104, 101, 108, 108, 111]); // 'hello'
const bytesValid = encoding.isValidUTF8Bytes(bytes);
console.log(bytesValid); // true

// Invalid UTF-8 bytes
const invalidBytes = new Uint8Array([0xFF, 0xFE]);
console.log(encoding.isValidUTF8Bytes(invalidBytes)); // false
```

#### Character and Byte Counting

```javascript
// Count bytes in UTF-8 string
const byteCount = encoding.countUTF8Bytes('Hello 🌍');
console.log(byteCount); // 10

// Count Unicode characters (runes)
const runeCount = encoding.countUTF8Runes('Hello 🌍');
console.log(runeCount); // 7

// Examples with different character types
console.log(encoding.countUTF8Bytes('你好')); // 6 bytes
console.log(encoding.countUTF8Runes('你好')); // 2 runes
console.log(encoding.countUTF8Bytes('café 🚀')); // 10 bytes
console.log(encoding.countUTF8Runes('café 🚀')); // 6 runes
```

#### Raw Bytes to String Conversion

```javascript
// Convert bytes to string without UTF-8 validation
const bytes = new Uint8Array([72, 101, 108, 108, 111]); // "Hello"
const str = encoding.bytesToString(bytes);
console.log(str); // "Hello"

// Handle binary data
const binary = new Uint8Array([0x00, 0xFF, 0x7F, 0x80]);
const binaryStr = encoding.bytesToString(binary);
console.log(binaryStr.length); // 4
```

### Error Handling

The extension provides proper error handling for invalid inputs:

```javascript
// Invalid UTF-8 bytes
try {
    encoding.decodeUTF8(new Uint8Array([0xFF, 0xFE]));
} catch (e) {
    console.log('Invalid UTF-8 bytes error:', e);
}

// Invalid Base64
try {
    encoding.decodeUTF8FromBase64('invalid base64!@#');
} catch (e) {
    console.log('Invalid Base64 error:', e);
}

// Null input
try {
    encoding.decodeUTF8(null);
} catch (e) {
    console.log('Null input error:', e);
}
```

### Performance Considerations

- The extension is optimized for high-performance text encoding/decoding
- UTF-8 validation is performed by default for safety
- Use `bytesToString` for raw byte conversion when UTF-8 validation is not needed
- Large inputs (up to 100MB) are supported
- Empty strings and null inputs are handled gracefully

### Supported Character Types

The extension handles various character types efficiently:
- ASCII characters
- Unicode characters
- Emoji and special symbols
- International characters (Chinese, Korean, Arabic, etc.)
- Mixed content (ASCII + Unicode + Emoji)

## License

This project is licensed under the MIT License - see the LICENSE file for details.