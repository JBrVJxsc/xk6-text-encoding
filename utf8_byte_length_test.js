import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  console.log('=== UTF-8 Byte Length Test ===');
  
  const textEncoding = new TextEncoding();
  
  // Test strings with various Unicode characters
  const testStrings = [
    "",
    "Hello",
    "Hello, 世界!",
    "🌍",
    "áéíóú",
    "1234567890",
    "!@#$%^&*()",
    "Hello, 世界! 🌍",
    "Special chars: áéíóú ñ ç ß € ¥ £",
    "Mixed: Hello 世界 🌍 áéíóú 123 !@#"
  ];
  
  console.log('\n1. Testing UTF8ByteLength function:');
  for (const str of testStrings) {
    const byteLength = textEncoding.utf8ByteLength(str);
    console.log(`"${str}": ${byteLength} bytes`);
  }
  
  console.log('\n2. Testing UTF8ByteLengthOptimized function:');
  for (const str of testStrings) {
    const byteLength = textEncoding.utf8ByteLengthOptimized(str);
    console.log(`"${str}": ${byteLength} bytes`);
  }
  
  console.log('\n3. Performance comparison with JavaScript implementation:');
  
  // JavaScript implementation for comparison
  function utf8ByteLengthJS(str) {
    let bytes = 0;
    for (let i = 0; i < str.length; i++) {
      const code = str.charCodeAt(i);
      if (code <= 0x7f) {
        bytes += 1;
      } else if (code <= 0x7ff) {
        bytes += 2;
      } else if (code >= 0xd800 && code <= 0xdbff) {
        // lead surrogate, assume a valid surrogate pair
        bytes += 4;
        i++; // skip low surrogate
      } else {
        bytes += 3;
      }
    }
    return bytes;
  }
  
  // Test consistency between Go and JavaScript implementations
  console.log('\n4. Consistency check between Go and JavaScript:');
  for (const str of testStrings) {
    const goResult = textEncoding.utf8ByteLength(str);
    const jsResult = utf8ByteLengthJS(str);
    
    if (goResult === jsResult) {
      console.log(`✓ "${str}": ${goResult} bytes (both implementations agree)`);
    } else {
      console.log(`✗ "${str}": Go=${goResult}, JS=${jsResult} (discrepancy!)`);
    }
  }
  
  console.log('\n5. Performance benchmark:');
  const longString = "Hello, 世界! 🌍 áéíóú ñ ç ß € ¥ £ ".repeat(1000);
  
  // Benchmark Go implementation
  const startGo = Date.now();
  for (let i = 0; i < 10000; i++) {
    textEncoding.utf8ByteLength(longString);
  }
  const endGo = Date.now();
  const goTime = endGo - startGo;
  
  // Benchmark JavaScript implementation
  const startJS = Date.now();
  for (let i = 0; i < 10000; i++) {
    utf8ByteLengthJS(longString);
  }
  const endJS = Date.now();
  const jsTime = endJS - startJS;
  
  console.log(`Go implementation: ${goTime}ms for 10,000 iterations`);
  console.log(`JS implementation: ${jsTime}ms for 10,000 iterations`);
  console.log(`Performance ratio: ${(jsTime / goTime).toFixed(2)}x faster in Go`);
  
  console.log('\n=== Test Complete ===');
} 