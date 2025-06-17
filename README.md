# xk6-text-encoding

A [k6](https://k6.io) extension providing UTF-8 text encoding and decoding capabilities for performance testing scenarios.

[![Go](https://github.com/JBrVJxsc/xk6-text-encoding/actions/workflows/go.yml/badge.svg)](https://github.com/JBrVJxsc/xk6-text-encoding/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- 🚀 **High Performance**: Optimized Go implementation for fast text encoding/decoding
- 🌍 **UTF-8 Support**: Full UTF-8 encoding and decoding capabilities
- 🔧 **Simple API**: Clean, modern JavaScript API
- ⚡ **Efficient UTF-8 Handling**: Fast UTF-8 byte and rune counting
- 🧪 **Well Tested**: Comprehensive test suite with benchmarks

## Installation

### Using xk6

```bash
# Build k6 with the text-encoding extension
xk6 build --with github.com/JBrVJxsc/xk6-text-encoding@latest
```

### Using Docker

```dockerfile
FROM grafana/xk6:latest as builder
RUN xk6 build --with github.com/JBrVJxsc/xk6-text-encoding@latest

FROM grafana/k6:latest
COPY --from=builder /xk6/k6 /usr/bin/k6
```

## Quick Start

```javascript
import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  const te = new TextEncoding();
  
  // Encoding text to UTF-8 bytes
  const bytes = te.EncodeUTF8('Hello, 世界! 🌍');
  
  // Encoding text to Base64
  const base64 = te.EncodeUTF8ToBase64('Hello, 世界! 🌍');
  
  // Decoding UTF-8 bytes back to text
  const decoded = te.DecodeUTF8(bytes);
  
  // Decoding Base64 back to text
  const decodedFromBase64 = te.DecodeUTF8FromBase64(base64);
  
  // Counting UTF-8 bytes
  const byteCount = te.CountUTF8Bytes('Hello, 世界! 🌍');
  
  // Counting UTF-8 runes (characters)
  const runeCount = te.CountUTF8Runes('Hello, 世界! 🌍');
  
  console.log(`Original: Hello, 世界! 🌍`);
  console.log(`Encoded bytes: ${bytes.length}`);
  console.log(`Base64: ${base64}`);
  console.log(`Decoded: ${decoded}`);
  console.log(`UTF-8 byte count: ${byteCount}`);
  console.log(`UTF-8 rune count: ${runeCount}`);
}
```

## API Reference

### TextEncoding

The main module that provides UTF-8 text encoding and decoding capabilities.

```javascript
import { TextEncoding } from 'k6/x/text-encoding';

const te = new TextEncoding();
```

#### Methods

##### `EncodeUTF8(text: string): Uint8Array`

Encodes a string to UTF-8 bytes.

```javascript
const te = new TextEncoding();
const bytes = te.EncodeUTF8('Hello World');
// Returns: Uint8Array containing the UTF-8 encoded bytes
```

##### `EncodeUTF8ToBase64(text: string): string`

Encodes a string to UTF-8 bytes and then to Base64.

```javascript
const te = new TextEncoding();
const base64 = te.EncodeUTF8ToBase64('Hello World');
// Returns: Base64 encoded string
```

##### `DecodeUTF8(data: Uint8Array): string`

Decodes UTF-8 bytes to a string. Throws an error if the bytes are not valid UTF-8.

```javascript
const te = new TextEncoding();
const bytes = new Uint8Array([72, 101, 108, 108, 111]); // "Hello"
const text = te.DecodeUTF8(bytes);
// Returns: "Hello"
```

##### `DecodeUTF8FromBase64(encodedData: string): string`

Decodes a Base64 string to UTF-8 text. Throws an error if the decoded data is not valid UTF-8.

```javascript
const te = new TextEncoding();
const base64 = 'SGVsbG8gV29ybGQ='; // "Hello World" in Base64
const text = te.DecodeUTF8FromBase64(base64);
// Returns: "Hello World"
```

##### `CountUTF8Bytes(text: string): number`

Returns the number of bytes in the UTF-8 encoding of the string.

```javascript
const te = new TextEncoding();
console.log(te.CountUTF8Bytes("Hello"));        // 5
console.log(te.CountUTF8Bytes("Hello, 世界!"));  // 14 
console.log(te.CountUTF8Bytes("🌍"));           // 4
```

##### `CountUTF8Runes(text: string): number`

Returns the number of Unicode code points (characters) in the string.

```javascript
const te = new TextEncoding();
console.log(te.CountUTF8Runes("Hello"));        // 5
console.log(te.CountUTF8Runes("Hello, 世界!"));  // 10
console.log(te.CountUTF8Runes("🌍"));           // 1
```

##### `IsValidUTF8(text: string): boolean`

Checks if a string is valid UTF-8.

```javascript
const te = new TextEncoding();
console.log(te.IsValidUTF8("Hello"));           // true
console.log(te.IsValidUTF8("Hello, 世界!"));     // true
console.log(te.IsValidUTF8("\uFFFD"));         // true
```

##### `IsValidUTF8Bytes(data: Uint8Array): boolean`

Checks if bytes represent valid UTF-8.

```javascript
const te = new TextEncoding();
const validBytes = new Uint8Array([72, 101, 108, 108, 111]); // "Hello"
const invalidBytes = new Uint8Array([0xFF, 0xFE]); // Invalid UTF-8
console.log(te.IsValidUTF8Bytes(validBytes));   // true
console.log(te.IsValidUTF8Bytes(invalidBytes)); // false
```

## Usage Examples

### Basic Encoding/Decoding

```javascript
import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  const te = new TextEncoding();
  
  const originalText = 'Hello, 世界! 🌍';
  const encoded = te.EncodeUTF8(originalText);
  const decoded = te.DecodeUTF8(encoded);
  
  console.log(`Roundtrip successful: ${originalText === decoded}`);
}
```

### Working with Base64

```javascript
import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  const te = new TextEncoding();
  
  const originalText = 'Hello, 世界! 🌍';
  const base64 = te.EncodeUTF8ToBase64(originalText);
  const decoded = te.DecodeUTF8FromBase64(base64);
  
  console.log(`Original: ${originalText}`);
  console.log(`Base64: ${base64}`);
  console.log(`Decoded: ${decoded}`);
}
```

### Character Counting

```javascript
import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  const te = new TextEncoding();
  
  const text = 'Hello, 世界! 🌍';
  const byteCount = te.CountUTF8Bytes(text);
  const runeCount = te.CountUTF8Runes(text);
  
  console.log(`Text: ${text}`);
  console.log(`Bytes: ${byteCount}`);
  console.log(`Characters: ${runeCount}`);
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [k6](https://k6.io) for the excellent performance testing platform
- [xk6](https://github.com/grafana/xk6) for the extension framework
- [golang.org/x/text](https://pkg.go.dev/golang.org/x/text) for text encoding support

## Support

- 📖 [Documentation](https://github.com/JBrVJxsc/xk6-text-encoding/blob/main/README.md)
- 🐛 [Issue Tracker](https://github.com/JBrVJxsc/xk6-text-encoding/issues)
- 💬 [Discussions](https://github.com/JBrVJxsc/xk6-text-encoding/discussions)
- 📧 [Email Support](mailto:your-email@example.com)

---

Made with ❤️ for the k6 community