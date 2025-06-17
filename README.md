# xk6-text-encoding

A [k6](https://k6.io) extension providing comprehensive text encoding and decoding capabilities for performance testing scenarios.

[![Go](https://github.com/JBrVJxsc/xk6-text-encoding/actions/workflows/go.yml/badge.svg)](https://github.com/JBrVJxsc/xk6-text-encoding/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- 🚀 **High Performance**: Optimized Go implementation for fast text encoding/decoding
- 🌍 **Multi-Language Support**: Supports 30+ text encodings including UTF-8, UTF-16, ISO-8859-*, Windows-*, Asian encodings
- 🔧 **Simple API**: Clean, modern JavaScript API similar to Web standards
- ⚡ **Efficient UTF-8 Handling**: Ultra-fast UTF-8 byte length calculation
- 🧪 **Well Tested**: Comprehensive test suite with benchmarks
- 📊 **Performance Monitoring**: Built-in performance comparison tools

## Supported Encodings

- **Unicode**: UTF-8, UTF-16 (LE/BE)
- **Western European**: ISO-8859-1 through ISO-8859-16, Windows-1250 through Windows-1258
- **Cyrillic**: KOI8-R, KOI8-U
- **Asian**: Shift-JIS, EUC-JP, ISO-2022-JP, GBK, GB18030, Big5, EUC-KR
- And many more...

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
import { TextEncoder, TextDecoder, Utils } from 'k6/x/text-encoding';

export default function () {
  // Encoding text
  const encoder = new TextEncoder('utf-8');
  const encoded = encoder.encode('Hello, 世界! 🌍');
  
  // Decoding data
  const decoder = new TextDecoder('utf-8');
  const decoded = decoder.decode(encoded);
  
  // Utility functions
  const utils = new Utils();
  const byteLength = utils.utf8ByteLength('Hello, 世界! 🌍');
  
  console.log(`Original: Hello, 世界! 🌍`);
  console.log(`Encoded bytes: ${encoded.length}`);
  console.log(`Decoded: ${decoded}`);
  console.log(`UTF-8 byte length: ${byteLength}`);
}
```

## API Reference

### TextEncoder

Creates an encoder for converting strings to bytes using a specified encoding.

```javascript
import { TextEncoder } from 'k6/x/text-encoding';

// Create encoder (defaults to UTF-8)
const encoder = new TextEncoder();
const utf8Encoder = new TextEncoder('utf-8');
const latin1Encoder = new TextEncoder('iso-8859-1');
```

#### Methods

##### `encode(text: string): Uint8Array`

Encodes a string to bytes using the specified encoding.

```javascript
const encoder = new TextEncoder('utf-8');
const bytes = encoder.encode('Hello World');
// Returns: Uint8Array containing the encoded bytes
```

##### `encodeString(text: string): string`

Convenience method that returns encoded bytes as a string.

```javascript
const encoder = new TextEncoder('utf-8');
const encodedString = encoder.encodeString('Hello World');
// Returns: String representation of encoded bytes
```

##### `getEncoding(): string`

Returns the encoding label used by this encoder.

```javascript
const encoder = new TextEncoder('utf-8');
console.log(encoder.getEncoding()); // "utf-8"
```

### TextDecoder

Creates a decoder for converting bytes to strings using a specified encoding.

```javascript
import { TextDecoder } from 'k6/x/text-encoding';

// Create decoder (defaults to UTF-8)
const decoder = new TextDecoder();
const utf8Decoder = new TextDecoder('utf-8');
const latin1Decoder = new TextDecoder('iso-8859-1');
```

#### Methods

##### `decode(data: Uint8Array): string`

Decodes bytes to a string using the specified encoding.

```javascript
const decoder = new TextDecoder('utf-8');
const bytes = new Uint8Array([72, 101, 108, 108, 111]); // "Hello"
const text = decoder.decode(bytes);
// Returns: "Hello"
```

##### `getEncoding(): string`

Returns the encoding label used by this decoder.

```javascript
const decoder = new TextDecoder('utf-8');
console.log(decoder.getEncoding()); // "utf-8"
```

### Utils

Provides utility functions for text encoding operations.

```javascript
import { Utils } from 'k6/x/text-encoding';

const utils = new Utils();
```

#### Methods

##### `utf8ByteLength(str: string): number`

Returns the byte length of a string when encoded in UTF-8. This is highly optimized and much faster than JavaScript alternatives.

```javascript
const utils = new Utils();
console.log(utils.utf8ByteLength("Hello"));        // 5
console.log(utils.utf8ByteLength("Hello, 世界!"));  // 14 
console.log(utils.utf8ByteLength("🌍"));           // 4
```

##### `isValidEncoding(label: string): boolean`

Checks if an encoding label is supported.

```javascript
const utils = new Utils();
console.log(utils.isValidEncoding("utf-8"));      // true
console.log(utils.isValidEncoding("invalid"));    // false
```

##### `getSupportedEncodings(): string[]`

Returns an array of all supported encoding labels.

```javascript
const utils = new Utils();
const encodings = utils.getSupportedEncodings();
console.log(`Supported encodings: ${encodings.length}`);
// Returns: ["utf-8", "utf-16", "iso-8859-1", ...]
```

## Usage Examples

### Basic Encoding/Decoding

```javascript
import { TextEncoder, TextDecoder } from 'k6/x/text-encoding';

export default function () {
  const encoder = new TextEncoder('utf-8');
  const decoder = new TextDecoder('utf-8');
  
  const originalText = 'Hello, 世界! 🌍';
  const encoded = encoder.encode(originalText);
  const decoded = decoder.decode(encoded);
  
  console.log(`Roundtrip successful: ${originalText === decoded}`);
}
```

### Working with Different Encodings

```javascript
import { TextEncoder, TextDecoder } from 'k6/x/text-encoding';

export default function () {
  const encodings = ['utf-8', 'utf-16', 'iso-8859-1', 'windows-1252'];
  const testText = 'Hello World!';
  
  encodings.forEach(encoding => {
    try {
      const encoder = new TextEncoder(encoding);
      const decoder = new TextDecoder(encoding);
      
      const encoded = encoder.encode(testText);
      const decoded = decoder.decode(encoded);
      
      console.log(`${encoding}: ${encoded.length} bytes`);
    } catch (error) {
      console.log(`${encoding}: Error - ${error.message}`);
    }
  });
}
```

### Performance Monitoring

```javascript
import { TextEncoder, Utils } from 'k6/x/text-encoding';

export default function () {
  const utils = new Utils();
  const encoder = new TextEncoder('utf-8');
  
  const testText = 'Performance test string with Unicode: 世界 🌍'.repeat(1000);
  
  // Measure UTF-8 byte length calculation
  const start1 = Date.now();
  const byteLength = utils.utf8ByteLength(testText);
  const end1 = Date.now();
  
  // Measure encoding
  const start2 = Date.now();
  const encoded = encoder.encode(testText);
  const end2 = Date.now();
  
  console.log(`UTF-8 byte length: ${byteLength} (${end1 - start1}ms)`);
  console.log(`Encoding: ${encoded.length} bytes (${end2 - start2}ms)`);
}
```

### Error Handling

```javascript
import { TextEncoder, Utils } from 'k6/x/text-encoding';

export default function () {
  const utils = new Utils();
  
  // Check encoding support before using
  const encoding = 'some-encoding';
  if (utils.isValidEncoding(encoding)) {
    const encoder = new TextEncoder(encoding);
    // Use encoder...
  } else {
    console.log(`Encoding ${encoding} is not supported`);
    
    // Show available alternatives
    const supported = utils.getSupportedEncodings();
    console.log(`Available encodings: ${supported.slice(0, 5).join(', ')}...`);
  }
}
```

## Building and Testing

### Prerequisites

- [Go](https://golang.org) 1.19+
- [xk6](https://github.com/grafana/xk6)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/JBrVJxsc/xk6-text-encoding.git
cd xk6-text-encoding

# Initialize the project
make init

# Run tests
make test-all

# Run benchmarks
make bench

# Compare Go vs JavaScript performance
make perf-compare
```

### Available Make Commands

#### Building
- `make build` - Build standalone extension
- `make build-xk6` - Build with xk6
- `make build-xk6-local` - Build with local code

#### Testing
- `make test-all` - Run both Go and K6 tests
- `make test-go` - Run Go tests only
- `make test-k6` - Run K6 tests only
- `make test-full` - Tests with coverage report

#### Benchmarking
- `make bench` - Run all Go benchmarks
- `make bench-utf8` - UTF-8 specific benchmarks
- `make perf-compare` - Compare Go vs K6 performance
- `make perf-full` - Complete performance test suite

#### Development
- `make dev` - Quick development cycle (format + lint + test)
- `make fmt` - Format code
- `make lint` - Lint code
- `make clean` - Clean build artifacts

## Performance

This extension is optimized for high performance:

- **UTF-8 byte length calculation**: Uses Go's native `len()` function, which is extremely fast for UTF-8 strings
- **Encoding/decoding**: Leverages Go's optimized `golang.org/x/text/encoding` package
- **Memory efficiency**: Built-in buffer pooling for memory reuse

### Benchmark Results

```
BenchmarkUTF8ByteLength-8           100000000    10.2 ns/op    0 B/op    0 allocs/op
BenchmarkUTF8ByteLengthJS-8          1000000   1023.5 ns/op    0 B/op    0 allocs/op
BenchmarkUTF8ByteLengthDirect-8     200000000     5.1 ns/op    0 B/op    0 allocs/op
```

The Go implementation is **~100x faster** than equivalent JavaScript implementations.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run the test suite (`make test-all`)
6. Run benchmarks (`make bench`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Development Guidelines

- Write tests for all new functionality
- Maintain backward compatibility
- Update documentation for API changes
- Run `make dev` before submitting PRs
- Add benchmarks for performance-critical code

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