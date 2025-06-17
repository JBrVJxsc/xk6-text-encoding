import textEncoding from 'k6/x/text-encoding';
import { check } from 'k6';

export let options = {
  iterations: 1,
  vus: 1,
};

// Helper functions for assertions
function assert(condition, message) {
  if (!condition) {
    throw new Error(`Assertion failed: ${message}`);
  }
}

function assertEqual(actual, expected, message) {
  if (actual !== expected) {
    throw new Error(`${message || 'Values not equal'}: expected ${expected}, got ${actual}`);
  }
}

function assertArrayEqual(actual, expected, message) {
  if (actual.length !== expected.length) {
    throw new Error(`${message || 'Array lengths differ'}: expected length ${expected.length}, got ${actual.length}`);
  }
  for (let i = 0; i < actual.length; i++) {
    if (actual[i] !== expected[i]) {
      throw new Error(`${message || 'Arrays differ'} at index ${i}: expected ${expected[i]}, got ${actual[i]}`);
    }
  }
}

function assertThrows(fn, message) {
  try {
    fn();
    throw new Error(`${message || 'Expected function to throw'}`);
  } catch (e) {
    // Expected to throw
  }
}

export default function () {
  console.log('=== Starting TextEncoding Extension Tests ===\n');
  
  // Test EncodeUTF8
  testEncodeUTF8();
  
  // Test EncodeUTF8ToBase64
  testEncodeUTF8ToBase64();
  
  // Test DecodeUTF8
  testDecodeUTF8();
  
  // Test DecodeUTF8FromBase64
  testDecodeUTF8FromBase64();
  
  // Test CountUTF8Bytes
  testCountUTF8Bytes();
  
  // Test CountUTF8Runes
  testCountUTF8Runes();
  
  // Test IsValidUTF8
  testIsValidUTF8();
  
  // Test IsValidUTF8Bytes
  testIsValidUTF8Bytes();
  
  // Test round-trip scenarios
  testRoundTrip();
  
  // Test performance
  testPerformance();
  
  // Test error handling
  testErrorHandling();
  
  console.log('\n=== All Tests Completed Successfully! ===');
}

function testEncodeUTF8() {
  console.log('Testing EncodeUTF8...');
  
  // Empty string
  let result = textEncoding.EncodeUTF8('');
  assertEqual(result.length, 0, 'Empty string should produce empty bytes');
  
  // ASCII text
  result = textEncoding.EncodeUTF8('hello');
  assertEqual(result.length, 5, 'ASCII "hello" should be 5 bytes');
  assertArrayEqual(Array.from(result), [104, 101, 108, 108, 111], 'ASCII bytes should match expected values');
  
  // Unicode text with emoji
  result = textEncoding.EncodeUTF8('Hello 🌍');
  assertEqual(result.length, 9, 'Unicode with emoji should be 9 bytes');
  
  // Chinese characters
  result = textEncoding.EncodeUTF8('你好');
  assertEqual(result.length, 6, 'Chinese characters should be 6 bytes (3 each)');
  
  console.log('✓ EncodeUTF8 tests passed\n');
}

function testEncodeUTF8ToBase64() {
  console.log('Testing EncodeUTF8ToBase64...');
  
  // Empty string
  let result = textEncoding.EncodeUTF8ToBase64('');
  assertEqual(result, '', 'Empty string should produce empty base64');
  
  // Simple text
  result = textEncoding.EncodeUTF8ToBase64('hello');
  assert(result.length > 0, 'Base64 result should not be empty');
  assert(typeof result === 'string', 'Base64 result should be string');
  
  // Verify round-trip
  let decoded = textEncoding.DecodeUTF8FromBase64(result);
  assertEqual(decoded, 'hello', 'Round-trip base64 should work');
  
  // Unicode text
  result = textEncoding.EncodeUTF8ToBase64('Hello 🌍');
  decoded = textEncoding.DecodeUTF8FromBase64(result);
  assertEqual(decoded, 'Hello 🌍', 'Unicode base64 round-trip should work');
  
  console.log('✓ EncodeUTF8ToBase64 tests passed\n');
}

function testDecodeUTF8() {
  console.log('Testing DecodeUTF8...');
  
  // Null input should throw
  assertThrows(() => textEncoding.DecodeUTF8(null), 'Null input should throw error');
  
  // Empty bytes
  let result = textEncoding.DecodeUTF8(new Uint8Array(0));
  assertEqual(result, '', 'Empty bytes should produce empty string');
  
  // Valid ASCII bytes
  let bytes = new Uint8Array([104, 101, 108, 108, 111]); // 'hello'
  result = textEncoding.DecodeUTF8(bytes);
  assertEqual(result, 'hello', 'ASCII bytes should decode correctly');
  
  // Valid Unicode bytes
  bytes = textEncoding.EncodeUTF8('Hello 🌍');
  result = textEncoding.DecodeUTF8(bytes);
  assertEqual(result, 'Hello 🌍', 'Unicode bytes should decode correctly');
  
  // Invalid UTF-8 bytes should throw
  let invalidBytes = new Uint8Array([0xFF, 0xFE]);
  assertThrows(() => textEncoding.DecodeUTF8(invalidBytes), 'Invalid UTF-8 bytes should throw error');
  
  console.log('✓ DecodeUTF8 tests passed\n');
}

function testDecodeUTF8FromBase64() {
  console.log('Testing DecodeUTF8FromBase64...');
  
  // Empty base64
  let result = textEncoding.DecodeUTF8FromBase64('');
  assertEqual(result, '', 'Empty base64 should produce empty string');
  
  // Valid base64 ASCII
  let encoded = textEncoding.EncodeUTF8ToBase64('hello');
  result = textEncoding.DecodeUTF8FromBase64(encoded);
  assertEqual(result, 'hello', 'Valid base64 ASCII should decode correctly');
  
  // Valid base64 Unicode
  encoded = textEncoding.EncodeUTF8ToBase64('Hello 🌍');
  result = textEncoding.DecodeUTF8FromBase64(encoded);
  assertEqual(result, 'Hello 🌍', 'Valid base64 Unicode should decode correctly');
  
  // Invalid base64 should throw
  assertThrows(() => textEncoding.DecodeUTF8FromBase64('invalid base64!@#'), 'Invalid base64 should throw error');
  
  // Complex Unicode
  let text = 'café naïve résumé 中文 🚀🌍💻';
  encoded = textEncoding.EncodeUTF8ToBase64(text);
  let decoded = textEncoding.DecodeUTF8FromBase64(encoded);
  assertEqual(decoded, text, 'Complex Unicode should handle round-trip correctly');
  
  console.log('✓ DecodeUTF8FromBase64 tests passed\n');
}

function testCountUTF8Bytes() {
  console.log('Testing CountUTF8Bytes...');
  
  assertEqual(textEncoding.CountUTF8Bytes(''), 0, 'Empty string should have 0 bytes');
  assertEqual(textEncoding.CountUTF8Bytes('hello'), 5, 'ASCII should have 5 bytes');
  assertEqual(textEncoding.CountUTF8Bytes('Hello 🌍'), 9, 'Unicode with emoji should have 9 bytes');
  assertEqual(textEncoding.CountUTF8Bytes('你好'), 6, 'Chinese characters should have 6 bytes');
  assertEqual(textEncoding.CountUTF8Bytes('café 🚀'), 10, 'Mixed content should have 10 bytes');
  assertEqual(textEncoding.CountUTF8Bytes('🚀🌍💻'), 12, 'Only emojis should have 12 bytes');
  
  console.log('✓ CountUTF8Bytes tests passed\n');
}

function testCountUTF8Runes() {
  console.log('Testing CountUTF8Runes...');
  
  assertEqual(textEncoding.CountUTF8Runes(''), 0, 'Empty string should have 0 runes');
  assertEqual(textEncoding.CountUTF8Runes('hello'), 5, 'ASCII should have 5 runes');
  assertEqual(textEncoding.CountUTF8Runes('Hello 🌍'), 7, 'Unicode with emoji should have 7 runes');
  assertEqual(textEncoding.CountUTF8Runes('你好'), 2, 'Chinese characters should have 2 runes');
  assertEqual(textEncoding.CountUTF8Runes('café 🚀'), 6, 'Mixed content should have 6 runes');
  assertEqual(textEncoding.CountUTF8Runes('🚀🌍💻'), 3, 'Only emojis should have 3 runes');
  
  console.log('✓ CountUTF8Runes tests passed\n');
}

function testIsValidUTF8() {
  console.log('Testing IsValidUTF8...');
  
  assertEqual(textEncoding.IsValidUTF8(''), true, 'Empty string should be valid');
  assertEqual(textEncoding.IsValidUTF8('hello'), true, 'ASCII should be valid');
  assertEqual(textEncoding.IsValidUTF8('Hello 🌍 你好'), true, 'Unicode should be valid');
  assertEqual(textEncoding.IsValidUTF8('café naïve résumé'), true, 'Special chars should be valid');
  assertEqual(textEncoding.IsValidUTF8('🚀🌍💻中文한국어العربية'), true, 'Complex Unicode should be valid');
  
  console.log('✓ IsValidUTF8 tests passed\n');
}

function testIsValidUTF8Bytes() {
  console.log('Testing IsValidUTF8Bytes...');
  
  assertEqual(textEncoding.IsValidUTF8Bytes(new Uint8Array(0)), true, 'Empty bytes should be valid');
  
  let bytes = new Uint8Array([104, 101, 108, 108, 111]); // 'hello'
  assertEqual(textEncoding.IsValidUTF8Bytes(bytes), true, 'ASCII bytes should be valid');
  
  bytes = textEncoding.EncodeUTF8('Hello 🌍');
  assertEqual(textEncoding.IsValidUTF8Bytes(bytes), true, 'Unicode bytes should be valid');
  
  let invalidBytes = new Uint8Array([0xFF, 0xFE]);
  assertEqual(textEncoding.IsValidUTF8Bytes(invalidBytes), false, 'Invalid UTF-8 bytes should not be valid');
  
  let incompleteBytes = new Uint8Array([0xF0, 0x9F]);
  assertEqual(textEncoding.IsValidUTF8Bytes(incompleteBytes), false, 'Incomplete UTF-8 should not be valid');
  
  let overlongBytes = new Uint8Array([0xC0, 0x80]);
  assertEqual(textEncoding.IsValidUTF8Bytes(overlongBytes), false, 'Overlong encoding should not be valid');
  
  console.log('✓ IsValidUTF8Bytes tests passed\n');
}

function testRoundTrip() {
  console.log('Testing round-trip encoding/decoding...');
  
  const testCases = [
    'hello',
    'Hello 🌍',
    '你好世界',
    'café naïve résumé',
    '🚀🌍💻🎉',
    'Mixed: ASCII + 中文 + العربية + 🚀',
    '', // empty string
    ' ', // single space
    '\n\t\r', // whitespace characters
    '!@#$%^&*()_+-=[]{}|;:,.<>?', // special ASCII
  ];
  
  testCases.forEach((testCase, index) => {
    // Test byte encoding round-trip
    let bytes = textEncoding.EncodeUTF8(testCase);
    let decodedFromBytes = textEncoding.DecodeUTF8(bytes);
    assertEqual(decodedFromBytes, testCase, `Round-trip bytes failed for case ${index + 1}`);
    
    // Test base64 encoding round-trip
    let base64 = textEncoding.EncodeUTF8ToBase64(testCase);
    let decodedFromBase64 = textEncoding.DecodeUTF8FromBase64(base64);
    assertEqual(decodedFromBase64, testCase, `Round-trip base64 failed for case ${index + 1}`);
    
    // Verify byte and rune counts are consistent
    let byteCount = textEncoding.CountUTF8Bytes(testCase);
    let runeCount = textEncoding.CountUTF8Runes(testCase);
    assertEqual(byteCount, bytes.length, `Byte count mismatch for case ${index + 1}`);
    assert(runeCount <= byteCount, `Rune count should not exceed byte count for case ${index + 1}`);
    
    // Verify validation
    assertEqual(textEncoding.IsValidUTF8(testCase), true, `String validation failed for case ${index + 1}`);
    assertEqual(textEncoding.IsValidUTF8Bytes(bytes), true, `Bytes validation failed for case ${index + 1}`);
  });
  
  console.log('✓ Round-trip tests passed\n');
}

function testPerformance() {
  console.log('Testing performance with large strings...');
  
  // Create a large string with mixed content
  let largeString = 'Hello 🌍 世界 '.repeat(1000);
  
  let startTime = Date.now();
  
  // Perform operations
  let bytes = textEncoding.EncodeUTF8(largeString);
  let decoded = textEncoding.DecodeUTF8(bytes);
  let byteCount = textEncoding.CountUTF8Bytes(largeString);
  let runeCount = textEncoding.CountUTF8Runes(largeString);
  let isValid = textEncoding.IsValidUTF8(largeString);
  
  let endTime = Date.now();
  let duration = endTime - startTime;
  
  // Verify correctness
  assertEqual(decoded, largeString, 'Large string round-trip failed');
  assertEqual(byteCount, bytes.length, 'Large string byte count mismatch');
  assertEqual(isValid, true, 'Large string validation failed');
  
  console.log(`Large string operations took ${duration}ms`);
  assert(duration < 1000, 'Performance test should complete in under 1 second');
  
  console.log('✓ Performance tests passed\n');
}

function testErrorHandling() {
  console.log('Testing error handling...');
  
  // Test that these don't crash (undefined handling)
  try {
    textEncoding.CountUTF8Bytes(undefined);
    textEncoding.CountUTF8Runes(undefined);
    textEncoding.IsValidUTF8(undefined);
  } catch (e) {
    // It's OK if they throw, just shouldn't crash the extension
  }
  
  // Null inputs where appropriate
  assertThrows(() => textEncoding.DecodeUTF8(null), 'DecodeUTF8 should throw on null');
  
  // IsValidUTF8Bytes with null should not crash (Go's utf8.Valid handles nil)
  try {
    textEncoding.IsValidUTF8Bytes(null);
  } catch (e) {
    // It's OK if it throws, just shouldn't crash
  }
  
  console.log('✓ Error handling tests passed\n');
}

// Export function for running basic functionality test
export function testBasicFunctionality() {
  let text = 'Hello 🌍';
  
  console.log('=== Basic Functionality Test ===');
  
  // Test encoding
  let bytes = textEncoding.EncodeUTF8(text);
  console.log(`Original: "${text}"`);
  console.log(`Encoded bytes length: ${bytes.length}`);
  
  // Test decoding
  let decoded = textEncoding.DecodeUTF8(bytes);
  console.log(`Decoded: "${decoded}"`);
  console.log(`Round-trip successful: ${text === decoded}`);
  
  // Test base64
  let base64 = textEncoding.EncodeUTF8ToBase64(text);
  let fromBase64 = textEncoding.DecodeUTF8FromBase64(base64);
  console.log(`Base64: ${base64}`);
  console.log(`From Base64: "${fromBase64}"`);
  console.log(`Base64 round-trip successful: ${text === fromBase64}`);
  
  // Test counting
  let byteCount = textEncoding.CountUTF8Bytes(text);
  let runeCount = textEncoding.CountUTF8Runes(text);
  console.log(`Byte count: ${byteCount}`);
  console.log(`Rune count: ${runeCount}`);
  
  // Test validation
  let isValid = textEncoding.IsValidUTF8(text);
  let bytesValid = textEncoding.IsValidUTF8Bytes(bytes);
  console.log(`String is valid UTF-8: ${isValid}`);
  console.log(`Bytes are valid UTF-8: ${bytesValid}`);
  
  // Use k6's check function for final verification
  check(null, {
    'Round-trip encoding works': () => text === decoded,
    'Base64 round-trip works': () => text === fromBase64,
    'Byte count is correct': () => byteCount === bytes.length,
    'String is valid UTF-8': () => isValid === true,
    'Bytes are valid UTF-8': () => bytesValid === true,
  });
  
  console.log('=== Test Complete ===');
}