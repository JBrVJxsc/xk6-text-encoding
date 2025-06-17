# xk6-text-encoding

A k6 extension that provides text encoding and decoding functionality using Go's built-in libraries. This extension supports various character encodings including UTF-8, UTF-16, ISO-8859 series, Windows code pages, and various Asian encodings.

## Installation

### Prerequisites

- Go 1.21 or later
- k6 v1.0.0 or later

### Building

```bash
# Clone the repository
git clone https://github.com/JBrVJxsc/xk6-text-encoding.git
cd xk6-text-encoding

# Build the extension
go build -o xk6-text-encoding .

# Or build with xk6
xk6 build --with github.com/JBrVJxsc/xk6-text-encoding@latest
```

## Usage

### Basic Usage

```javascript
import { TextEncoder, TextDecoder } from 'k6/x/text-encoding';

export default function () {
  // Create encoder and decoder instances
  const encoder = new TextEncoder("utf-8");
  const decoder = new TextDecoder("utf-8");
  
  // Encode text to bytes
  const text = "Hello, 世界!";
  const encoded = encoder.encode(text);
  
  // Decode bytes back to text
  const decoded = decoder.decode(encoded);
  
  console.log(`Original: ${text}`);
  console.log(`Encoded: ${encoded}`);
  console.log(`Decoded: ${decoded}`);
}
```

### UTF-8 Byte Length

```javascript
import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  const textEncoding = new TextEncoding();
  
  // Get UTF-8 byte length of strings
  const strings = ["Hello", "Hello, 世界!", "🌍", "áéíóú"];
  
  for (const str of strings) {
    const byteLength = textEncoding.utf8ByteLength(str);
    console.log(`"${str}": ${byteLength} bytes`);
  }
  
  // Performance comparison with JavaScript implementation
  const longString = "Hello, 世界! 🌍 ".repeat(1000);
  
  // Much faster than JavaScript implementation
  const start = Date.now();
  for (let i = 0; i < 10000; i++) {
    textEncoding.utf8ByteLength(longString);
  }
  const end = Date.now();
  
  console.log(`Processed 10,000 iterations in ${end - start}ms`);
}
```

### Different Encodings

```javascript
import { TextEncoder, TextDecoder } from 'k6/x/text-encoding';

export default function () {
  // UTF-8 encoding (default)
  const utf8Encoder = new TextEncoder("utf-8");
  const utf8Decoder = new TextDecoder("utf-8");
  
  // ISO-8859-1 (Latin-1) encoding
  const latin1Encoder = new TextEncoder("iso-8859-1");
  const latin1Decoder = new TextDecoder("iso-8859-1");
  
  // Shift-JIS encoding for Japanese
  const sjisEncoder = new TextEncoder("shift-jis");
  const sjisDecoder = new TextDecoder("shift-jis");
  
  const text = "Hello, 世界!";
  
  // Encode with different encodings
  const utf8Bytes = utf8Encoder.encode(text);
  const latin1Bytes = latin1Encoder.encode(text);
  const sjisBytes = sjisEncoder.encode(text);
  
  console.log(`UTF-8 bytes: ${utf8Bytes.length}`);
  console.log(`Latin-1 bytes: ${latin1Bytes.length}`);
  console.log(`Shift-JIS bytes: ${sjisBytes.length}`);
}
```

### Error Handling

```javascript
import { TextEncoder, TextDecoder } from 'k6/x/text-encoding';

export default function () {
  try {
    // Try to create an encoder with unsupported encoding
    const encoder = new TextEncoder("unsupported-encoding");
  } catch (error) {
    console.log(`Error: ${error.message}`);
  }
  
  try {
    // Try to decode invalid data
    const decoder = new TextDecoder("utf-8");
    const decoded = decoder.decode([0xFF, 0xFE, 0x00]); // Invalid UTF-8
    console.log(`Decoded: ${decoded}`);
  } catch (error) {
    console.log(`Decode error: ${error.message}`);
  }
}
```

### Utility Functions

```javascript
import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  const textEncoding = new TextEncoding();
  
  // Check if encoding is supported
  console.log(`UTF-8 supported: ${textEncoding.isValidEncoding("utf-8")}`);
  console.log(`Invalid encoding supported: ${textEncoding.isValidEncoding("invalid")}`);
  
  // Get list of supported encodings
  const supported = textEncoding.getSupportedEncodings();
  console.log(`Supported encodings: ${supported.join(", ")}`);
}
```

## API Reference

### TextEncoder

#### Constructor
```javascript
const encoder = new TextEncoder(encoding);
```
- `encoding` (string, optional): The encoding to use. Defaults to "utf-8".

#### Methods

##### `encode(text)`
Encodes a string to bytes using the specified encoding.
- `text` (string): The text to encode.
- Returns: `[]byte` - The encoded bytes.

##### `encodeString(text)`
Convenience method that returns the encoded bytes as a string.
- `text` (string): The text to encode.
- Returns: `string` - The encoded bytes as a string.

##### `getEncoding()`
Returns the encoding label used by this encoder.
- Returns: `string` - The encoding label.

### TextDecoder

#### Constructor
```javascript
const decoder = new TextDecoder(encoding);
```
- `encoding` (string, optional): The encoding to use. Defaults to "utf-8".

#### Methods

##### `decode(data)`
Decodes bytes to a string using the specified encoding.
- `data` ([]byte): The bytes to decode.
- Returns: `string` - The decoded text.

##### `getEncoding()`
Returns the encoding label used by this decoder.
- Returns: `string` - The encoding label.

### TextEncoding (Utility Class)

#### Methods

##### `isValidEncoding(label)`
Checks if the given encoding label is supported.
- `label` (string): The encoding label to check.
- Returns: `boolean` - True if the encoding is supported.

##### `getSupportedEncodings()`
Returns a list of all supported encoding labels.
- Returns: `[]string` - Array of supported encoding labels.

##### `utf8ByteLength(str)`
Returns the byte length of a string when encoded in UTF-8. This is much faster than JavaScript implementations.
- `str` (string): The input string to measure.
- Returns: `number` - Number of bytes when encoded in UTF-8.

##### `utf8ByteLengthOptimized(str)`
Alternative implementation that manually calculates UTF-8 byte length. Useful for educational purposes.
- `str` (string): The input string to measure.
- Returns: `number` - Number of bytes when encoded in UTF-8.

## Supported Encodings

### Unicode
- `utf-8`, `utf8`
- `utf-16`, `utf16`
- `utf-16le`, `utf16le`
- `utf-16be`, `utf16be`

### ISO-8859 Series
- `iso-8859-1`, `latin1`
- `iso-8859-2`, `latin2`
- `iso-8859-3`, `latin3`
- `iso-8859-4`, `latin4`
- `iso-8859-5`
- `iso-8859-6`
- `iso-8859-7`
- `iso-8859-8`
- `iso-8859-9`, `latin5`
- `iso-8859-10`, `latin6`
- `iso-8859-13`, `latin7`
- `iso-8859-14`, `latin8`
- `iso-8859-15`, `latin9`
- `iso-8859-16`, `latin10`

### Windows Code Pages
- `windows-1250`
- `windows-1251`
- `windows-1252`
- `windows-1253`
- `windows-1254`
- `windows-1255`
- `windows-1256`
- `windows-1257`
- `windows-1258`

### Cyrillic
- `koi8-r`
- `koi8-u`

### Japanese
- `shift-jis`, `shift_jis`, `sjis`
- `euc-jp`, `eucjp`
- `iso-2022-jp`, `iso2022jp`

### Chinese
- `gbk`
- `gb18030`
- `big5`

### Korean
- `euc-kr`, `euckr`
- `iso-2022-kr`, `iso2022kr`

## Performance Considerations

- The extension uses buffer pooling to reduce memory allocations
- UTF-8 encoding/decoding is optimized for direct byte conversion
- For other encodings, the Go `golang.org/x/text/encoding` package is used

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
