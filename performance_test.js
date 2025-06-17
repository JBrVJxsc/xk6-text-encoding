import { TextEncoder, TextDecoder, Utils } from 'k6/x/text-encoding';
import { check } from 'k6';

/**
 * Performance test comparing three UTF-8 byte length implementations:
 * 1. Go Utils.utf8ByteLength() - optimized native implementation
 * 2. JavaScript loop-based calculation - manual UTF-8 byte counting
 * 3. TextEncoder.encode().length - using our extension's encoder
 */

// Test data - various string types for comprehensive benchmarking
const testStrings = {
  ascii: "Hello World! This is a simple ASCII string for testing performance.",
  unicode: "Hello, 世界! 🌍 This contains Unicode: áéíóú ñ ç ß € ¥ £",
  emoji: "🌍🌎🌏🚀⭐🌟💫✨🔥💧🌊🌈☀️🌙",
  mixed: "ASCII mixed with 世界 emoji 🌍 and special chars: áéíóú",
  large: "Hello, 世界! 🌍 ".repeat(100), // About 1.6KB
  huge: "Mixed content with ASCII, 世界, 🌍, and áéíóú. ".repeat(1000) // About 50KB
};

// Create encoder instances once for reuse
const utils = new Utils();
const encoder = new TextEncoder('utf-8');

// JavaScript implementation for comparison (loop-based)
function utf8ByteLengthJS(str) {
  let bytes = 0;
  for (let i = 0; i < str.length; i++) {
    const code = str.charCodeAt(i);
    if (code <= 0x7F) {
      bytes += 1;
    } else if (code <= 0x7FF) {
      bytes += 2;
    } else if ((code & 0xFC00) === 0xD800 && i + 1 < str.length) {
      // High surrogate
      const low = str.charCodeAt(++i);
      if ((low & 0xFC00) === 0xDC00) {
        // Valid surrogate pair - 4 bytes
        bytes += 4;
      } else {
        // Invalid surrogate
        bytes += 3; // Treat as replacement character
        i--; // Back up one
      }
    } else {
      bytes += 3;
    }
  }
  return bytes;
}

// Using our TextEncoder for comparison
function utf8ByteLengthViaEncoder(str) {
  return encoder.encode(str).length;
}

export const options = {
  vus: 1,
  duration: '10s',
  thresholds: {
    checks: ['rate>0.95']
  }
};

export default function() {
  let checksPass = 0;
  let totalChecks = 0;
  
  // Performance counters
  let goTotalTime = 0;
  let jsTotalTime = 0;
  let encTotalTime = 0;
  let iterations = 0;
  
  // Test each string type
  for (const [testType, testString] of Object.entries(testStrings)) {
    iterations++;
    
    // Test 1: Go implementation (Utils.utf8ByteLength)
    const startGo = Date.now();
    const goResult = utils.utf8ByteLength(testString);
    const endGo = Date.now();
    goTotalTime += (endGo - startGo);
    
    // Test 2: JavaScript loop implementation
    const startJS = Date.now();
    const jsResult = utf8ByteLengthJS(testString);
    const endJS = Date.now();
    jsTotalTime += (endJS - startJS);
    
    // Test 3: TextEncoder implementation (via encoding)
    const startEnc = Date.now();
    const encResult = utf8ByteLengthViaEncoder(testString);
    const endEnc = Date.now();
    encTotalTime += (endEnc - startEnc);
    
    // Verify all implementations give the same result
    totalChecks += 3;
    if (goResult === jsResult) checksPass++;
    if (goResult === encResult) checksPass++;
    if (jsResult === encResult) checksPass++;
    
    check(null, {
      [`${testType}: Go vs JS match`]: () => goResult === jsResult,
      [`${testType}: Go vs Encoder match`]: () => goResult === encResult,
      [`${testType}: JS vs Encoder match`]: () => jsResult === encResult,
    });
    
    // Log detailed results for this iteration
    if (__ITER % 100 === 0) { // Log every 100th iteration to avoid spam
      console.log(`${testType} (${testString.length} chars):`);
      console.log(`  Go Utils: ${goResult} bytes (${endGo - startGo}ms)`);
      console.log(`  JS Loop: ${jsResult} bytes (${endJS - startJS}ms)`);
      console.log(`  Encoder: ${encResult} bytes (${endEnc - startEnc}ms)`);
    }
  }
  
  // Log performance summary every 50 iterations
  if (__ITER % 50 === 0 && __ITER > 0) {
    const avgGoTime = goTotalTime / iterations;
    const avgJSTime = jsTotalTime / iterations;
    const avgEncTime = encTotalTime / iterations;
    
    console.log(`\n=== Performance Summary (Iteration ${__ITER}) ===`);
    console.log(`Average time per test set:`);
    console.log(`  Go Utils: ${avgGoTime.toFixed(4)}ms`);
    console.log(`  JS Loop: ${avgJSTime.toFixed(4)}ms`);
    console.log(`  TextEncoder: ${avgEncTime.toFixed(4)}ms`);
    
    if (avgJSTime > 0) {
      const speedupVsJS = (avgJSTime / avgGoTime).toFixed(2);
      console.log(`  Go is ${speedupVsJS}x faster than JS loop`);
    }
    
    if (avgEncTime > 0) {
      const speedupVsEnc = (avgEncTime / avgGoTime).toFixed(2);
      console.log(`  Go is ${speedupVsEnc}x faster than TextEncoder`);
    }
    
    const accuracy = ((checksPass / totalChecks) * 100).toFixed(2);
    console.log(`  Accuracy: ${accuracy}% (${checksPass}/${totalChecks})`);
    console.log('=======================================\n');
  }
  
  // Test encoding/decoding performance
  if (__ITER % 20 === 0) {
    const decoder = new TextDecoder('utf-8');
    const testText = testStrings.mixed;
    
    const encodeStart = Date.now();
    const encoded = encoder.encode(testText);
    const encodeEnd = Date.now();
    
    const decodeStart = Date.now();
    const decoded = decoder.decode(encoded);
    const decodeEnd = Date.now();
    
    check(null, {
      'Encoding/Decoding roundtrip': () => decoded === testText,
    });
    
    if (__ITER % 100 === 0) {
      console.log(`Encoding/Decoding performance:`);
      console.log(`  Encode: ${encodeEnd - encodeStart}ms`);
      console.log(`  Decode: ${decodeEnd - decodeStart}ms`);
      console.log(`  Roundtrip success: ${decoded === testText}`);
    }
  }
  
  // Stress test with huge strings occasionally
  if (__ITER % 100 === 0) {
    const hugeString = testStrings.huge;
    const start = Date.now();
    const result = utils.utf8ByteLength(hugeString);
    const end = Date.now();
    
    console.log(`Stress test (${hugeString.length} chars): ${result} bytes in ${end - start}ms`);
  }
}

export function teardown() {
  console.log('\n🏁 Performance test completed!');
  console.log('📊 Check the logs above for detailed performance comparisons');
  console.log('🚀 Comparison: Go Utils vs JavaScript loop vs TextEncoder');
  console.log('💡 Go implementation should show significant performance advantages');
}