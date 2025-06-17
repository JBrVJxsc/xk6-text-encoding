# xk6-text-encoding

A k6 extension for text encoding and decoding operations.

## Features

- UTF-8 encoding and decoding
- Base64 encoding and decoding
- UTF-8 validation
- Byte and rune counting
- High performance
- Simple API

## Installation

```bash
go install go.k6.io/xk6/cmd/xk6@latest
xk6 build --with github.com/JBrVJxsc/xk6-text-encoding@latest
```

## Quick Start

```javascript
import encoding from 'k6/x/text-encoding';

export default function () {
  // Create a new text encoder
  const encoder = new encoding.TextEncoder();
  
  // Create a new text decoder
  const decoder = new encoding.TextDecoder();

  // Encode a string to UTF-8 bytes
  const bytes = encoder.encode('Hello 🌍');
  console.log('Encoded bytes:', bytes);

  // Decode UTF-8 bytes back to string
  const text = decoder.decode(bytes);
  console.log('Decoded text:', text);

  // Encode to Base64
  const base64 = encoder.encodeToBase64('Hello 🌍');
  console.log('Base64:', base64);

  // Decode from Base64
  const decoded = decoder.decodeFromBase64(base64);
  console.log('Decoded from Base64:', decoded);

  // Count bytes and runes
  const byteCount = encoder.countBytes('Hello 🌍');
  console.log('Byte count:', byteCount);

  const runeCount = encoder.countRunes('Hello 🌍');
  console.log('Rune count:', runeCount);

  // Validate UTF-8
  const isValid = encoder.isValid('Hello 🌍');
  console.log('Is valid UTF-8:', isValid);
}
```

## API Reference

### TextEncoder

The `TextEncoder` class provides methods for encoding text and performing UTF-8 operations.

#### Methods

- `encode(text: string): Uint8Array` - Encodes a string to UTF-8 bytes
- `encodeToBase64(text: string): string` - Encodes a string to Base64
- `countBytes(text: string): number` - Counts the number of bytes in a UTF-8 string
- `countRunes(text: string): number` - Counts the number of runes (Unicode characters) in a string
- `isValid(text: string): boolean` - Checks if a string is valid UTF-8

### TextDecoder

The `TextDecoder` class provides methods for decoding text from various encodings.

#### Methods

- `decode(bytes: Uint8Array): string` - Decodes UTF-8 bytes to a string
- `decodeFromBase64(base64: string): string` - Decodes a Base64 string to a UTF-8 string

## Examples

### Basic Encoding and Decoding

```javascript
import encoding from 'k6/x/text-encoding';

export default function () {
  const encoder = new encoding.TextEncoder();
  const decoder = new encoding.TextDecoder();

  // Simple ASCII text
  const asciiText = 'Hello, World!';
  const asciiBytes = encoder.encode(asciiText);
  console.log('ASCII bytes:', asciiBytes);
  console.log('Decoded ASCII:', decoder.decode(asciiBytes));

  // Unicode text with emoji
  const unicodeText = 'Hello 🌍 你好 안녕하세요';
  const unicodeBytes = encoder.encode(unicodeText);
  console.log('Unicode bytes:', unicodeBytes);
  console.log('Decoded Unicode:', decoder.decode(unicodeBytes));
}
```

### Base64 Operations

```javascript
import encoding from 'k6/x/text-encoding';

export default function () {
  const encoder = new encoding.TextEncoder();
  const decoder = new encoding.TextDecoder();

  // Encode to Base64
  const text = 'Hello 🌍 你好';
  const base64 = encoder.encodeToBase64(text);
  console.log('Base64:', base64);

  // Decode from Base64
  const decoded = decoder.decodeFromBase64(base64);
  console.log('Decoded:', decoded);
}
```

### UTF-8 Validation and Counting

```javascript
import encoding from 'k6/x/text-encoding';

export default function () {
  const encoder = new encoding.TextEncoder();

  const text = 'Hello 🌍 你好 안녕하세요';

  // Count bytes and runes
  const byteCount = encoder.countBytes(text);
  const runeCount = encoder.countRunes(text);
  console.log(`Text has ${byteCount} bytes and ${runeCount} runes`);

  // Validate UTF-8
  const isValid = encoder.isValid(text);
  console.log('Is valid UTF-8:', isValid);
}
```

### Error Handling

```javascript
import encoding from 'k6/x/text-encoding';

export default function () {
  const encoder = new encoding.TextEncoder();
  const decoder = new encoding.TextDecoder();

  try {
    // Try to decode invalid UTF-8 bytes
    const invalidBytes = new Uint8Array([0xFF, 0xFE]);
    decoder.decode(invalidBytes);
  } catch (e) {
    console.log('Error:', e.message);
  }

  try {
    // Try to decode invalid Base64
    decoder.decodeFromBase64('invalid base64!');
  } catch (e) {
    console.log('Error:', e.message);
  }
}
```

## Performance Considerations

- The extension is optimized for high performance
- UTF-8 validation is done efficiently using Go's built-in UTF-8 validation
- Base64 operations use Go's standard library implementation
- Memory allocations are minimized where possible

## License

This project is licensed under the MIT License - see the LICENSE file for details.