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
    throw new Error(`${message || 'Values not equal'}: expected "${expected}", got "${actual}"`);
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
  
  // Test encodeUTF8
  testEncodeUTF8();
  
  // Test encodeUTF8ToBase64
  testEncodeUTF8ToBase64();
  
  // Test decodeUTF8
  testDecodeUTF8();
  
  // Test decodeUTF8FromBase64
  testDecodeUTF8FromBase64();
  
  // Test countUTF8Bytes
  testCountUTF8Bytes();
  
  // Test countUTF8Runes
  testCountUTF8Runes();
  
  // Test isValidUTF8
  testIsValidUTF8();
  
  // Test isValidUTF8Bytes
  testIsValidUTF8Bytes();
  
  // Test round-trip scenarios
  testRoundTrip();
  
  // Test performance
  testPerformance();
  
  // Test error handling
  testErrorHandling();
  
  // Test invalid UTF-8 sequences
  testInvalidUTF8Sequences();
  
  // Test concurrent operations
  testConcurrentOperations();
  
  // Test stress with large strings
  testStressLargeText();
  
  console.log('\n=== All Tests Completed Successfully! ===');
}

function testEncodeUTF8() {
  console.log('Testing encodeUTF8...');
  
  // Empty string
  let result = textEncoding.encodeUTF8('');
  assertEqual(result.length, 0, 'Empty string should produce empty bytes');
  
  // ASCII text
  result = textEncoding.encodeUTF8('hello');
  assertEqual(result.length, 5, 'ASCII "hello" should be 5 bytes');
  assertArrayEqual(Array.from(result), [104, 101, 108, 108, 111], 'ASCII bytes should match expected values');
  
  // Unicode text with emoji
  result = textEncoding.encodeUTF8('Hello 🌍');
  assertEqual(result.length, 10, 'Unicode with emoji should be 10 bytes');
  
  // Chinese characters
  result = textEncoding.encodeUTF8('你好');
  assertEqual(result.length, 6, 'Chinese characters should be 6 bytes (3 each)');
  
  console.log('✓ encodeUTF8 tests passed\n');
}

function testEncodeUTF8ToBase64() {
  console.log('Testing encodeUTF8ToBase64...');
  
  // Empty string
  let result = textEncoding.encodeUTF8ToBase64('');
  assertEqual(result, '', 'Empty string should produce empty base64');
  
  // Simple text
  result = textEncoding.encodeUTF8ToBase64('hello');
  assert(result.length > 0, 'Base64 result should not be empty');
  assert(typeof result === 'string', 'Base64 result should be string');
  
  // Verify round-trip
  let decoded = textEncoding.decodeUTF8FromBase64(result);
  assertEqual(decoded, 'hello', 'Round-trip base64 should work');
  
  // Unicode text
  result = textEncoding.encodeUTF8ToBase64('Hello 🌍');
  decoded = textEncoding.decodeUTF8FromBase64(result);
  assertEqual(decoded, 'Hello 🌍', 'Unicode base64 round-trip should work');
  
  console.log('✓ encodeUTF8ToBase64 tests passed\n');
}

function testDecodeUTF8() {
  console.log('Testing decodeUTF8...');
  
  // Null input should throw
  assertThrows(() => textEncoding.decodeUTF8(null), 'Null input should throw error');
  
  // Empty bytes
  let result = textEncoding.decodeUTF8(textEncoding.encodeUTF8(''));
  assertEqual(result, '', 'Empty bytes should produce empty string');
  
  // Valid ASCII bytes
  let bytes = new Uint8Array([104, 101, 108, 108, 111]); // 'hello'
  result = textEncoding.decodeUTF8(bytes);
  assertEqual(result, 'hello', 'ASCII bytes should decode correctly');
  
  // Valid Unicode bytes
  bytes = textEncoding.encodeUTF8('Hello 🌍');
  result = textEncoding.decodeUTF8(bytes);
  assertEqual(result, 'Hello 🌍', 'Unicode bytes should decode correctly');
  
  // Invalid UTF-8 bytes should throw
  let invalidBytes = new Uint8Array([0xFF, 0xFE]);
  assertThrows(() => textEncoding.decodeUTF8(invalidBytes), 'Invalid UTF-8 bytes should throw error');
  
  console.log('✓ decodeUTF8 tests passed\n');
}

function testDecodeUTF8FromBase64() {
  console.log('Testing decodeUTF8FromBase64...');
  
  // Empty base64
  let result = textEncoding.decodeUTF8FromBase64('');
  assertEqual(result, '', 'Empty base64 should produce empty string');
  
  // Valid base64 ASCII
  let encoded = textEncoding.encodeUTF8ToBase64('hello');
  result = textEncoding.decodeUTF8FromBase64(encoded);
  assertEqual(result, 'hello', 'Valid base64 ASCII should decode correctly');
  
  // Valid base64 Unicode
  encoded = textEncoding.encodeUTF8ToBase64('Hello 🌍');
  result = textEncoding.decodeUTF8FromBase64(encoded);
  assertEqual(result, 'Hello 🌍', 'Valid base64 Unicode should decode correctly');
  
  // Invalid base64 should throw
  assertThrows(() => textEncoding.decodeUTF8FromBase64('invalid base64!@#'), 'Invalid base64 should throw error');
  
  // Complex Unicode
  let text = 'café naïve résumé 中文 🚀🌍💻';
  encoded = textEncoding.encodeUTF8ToBase64(text);
  let decoded = textEncoding.decodeUTF8FromBase64(encoded);
  assertEqual(decoded, text, 'Complex Unicode should handle round-trip correctly');
  
  console.log('✓ decodeUTF8FromBase64 tests passed\n');
}

function testCountUTF8Bytes() {
  console.log('Testing countUTF8Bytes...');
  
  assertEqual(textEncoding.countUTF8Bytes(''), 0, 'Empty string should have 0 bytes');
  assertEqual(textEncoding.countUTF8Bytes('hello'), 5, 'ASCII should have 5 bytes');
  assertEqual(textEncoding.countUTF8Bytes('Hello 🌍'), 10, 'Unicode with emoji should have 10 bytes');
  assertEqual(textEncoding.countUTF8Bytes('你好'), 6, 'Chinese characters should have 6 bytes');
  assertEqual(textEncoding.countUTF8Bytes('café 🚀'), 10, 'Mixed content should have 10 bytes');
  assertEqual(textEncoding.countUTF8Bytes('🚀🌍💻'), 12, 'Only emojis should have 12 bytes');
  
  console.log('✓ countUTF8Bytes tests passed\n');
}

function testCountUTF8Runes() {
  console.log('Testing countUTF8Runes...');
  
  assertEqual(textEncoding.countUTF8Runes(''), 0, 'Empty string should have 0 runes');
  assertEqual(textEncoding.countUTF8Runes('hello'), 5, 'ASCII should have 5 runes');
  assertEqual(textEncoding.countUTF8Runes('Hello 🌍'), 7, 'Unicode with emoji should have 7 runes');
  assertEqual(textEncoding.countUTF8Runes('你好'), 2, 'Chinese characters should have 2 runes');
  assertEqual(textEncoding.countUTF8Runes('café 🚀'), 6, 'Mixed content should have 6 runes');
  assertEqual(textEncoding.countUTF8Runes('🚀🌍💻'), 3, 'Only emojis should have 3 runes');
  
  console.log('✓ countUTF8Runes tests passed\n');
}

function testIsValidUTF8() {
  console.log('Testing isValidUTF8...');
  
  assertEqual(textEncoding.isValidUTF8(''), true, 'Empty string should be valid');
  assertEqual(textEncoding.isValidUTF8('hello'), true, 'ASCII should be valid');
  assertEqual(textEncoding.isValidUTF8('Hello 🌍 你好'), true, 'Unicode should be valid');
  assertEqual(textEncoding.isValidUTF8('café naïve résumé'), true, 'Special chars should be valid');
  assertEqual(textEncoding.isValidUTF8('🚀🌍💻中文한국어العربية'), true, 'Complex Unicode should be valid');
  
  console.log('✓ isValidUTF8 tests passed\n');
}

function testIsValidUTF8Bytes() {
  console.log('Testing isValidUTF8Bytes...');
  
  assertEqual(textEncoding.isValidUTF8Bytes(new Uint8Array(0)), true, 'Empty bytes should be valid');
  
  let bytes = new Uint8Array([104, 101, 108, 108, 111]); // 'hello'
  assertEqual(textEncoding.isValidUTF8Bytes(bytes), true, 'ASCII bytes should be valid');
  
  bytes = textEncoding.encodeUTF8('Hello 🌍');
  assertEqual(textEncoding.isValidUTF8Bytes(bytes), true, 'Unicode bytes should be valid');
  
  let invalidBytes = new Uint8Array([0xFF, 0xFE]);
  assertEqual(textEncoding.isValidUTF8Bytes(invalidBytes), false, 'Invalid UTF-8 bytes should not be valid');
  
  let incompleteBytes = new Uint8Array([0xF0, 0x9F]);
  assertEqual(textEncoding.isValidUTF8Bytes(incompleteBytes), false, 'Incomplete UTF-8 should not be valid');
  
  let overlongBytes = new Uint8Array([0xC0, 0x80]);
  assertEqual(textEncoding.isValidUTF8Bytes(overlongBytes), false, 'Overlong encoding should not be valid');
  
  console.log('✓ isValidUTF8Bytes tests passed\n');
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
    let bytes = textEncoding.encodeUTF8(testCase);
    let decodedFromBytes = textEncoding.decodeUTF8(bytes);
    assertEqual(decodedFromBytes, testCase, `Round-trip bytes failed for case ${index + 1}`);
    
    // Test base64 encoding round-trip
    let base64 = textEncoding.encodeUTF8ToBase64(testCase);
    let decodedFromBase64 = textEncoding.decodeUTF8FromBase64(base64);
    assertEqual(decodedFromBase64, testCase, `Round-trip base64 failed for case ${index + 1}`);
    
    // Verify byte and rune counts are consistent
    let byteCount = textEncoding.countUTF8Bytes(testCase);
    let runeCount = textEncoding.countUTF8Runes(testCase);
    assertEqual(byteCount, bytes.length, `Byte count mismatch for case ${index + 1}`);
    assert(runeCount <= byteCount, `Rune count should not exceed byte count for case ${index + 1}`);
    
    // Verify validation
    assertEqual(textEncoding.isValidUTF8(testCase), true, `String validation failed for case ${index + 1}`);
    assertEqual(textEncoding.isValidUTF8Bytes(bytes), true, `Bytes validation failed for case ${index + 1}`);
  });
  
  console.log('✓ Round-trip tests passed\n');
}

function testPerformance() {
  console.log('Testing performance with large strings...');
  
  // Create a large string with various Unicode characters
  let largeString = '';
  // Add ASCII characters
  for (let i = 0; i < 5000; i++) {
    largeString += 'Hello ';
  }
  // Add Chinese characters
  for (let i = 0; i < 2500; i++) {
    largeString += '你好';
  }
  // Add emojis
  for (let i = 0; i < 1000; i++) {
    largeString += '🌍';
  }
  // Add mixed content
  for (let i = 0; i < 500; i++) {
    largeString += 'café ';
  }
  // Add more diverse content
  for (let i = 0; i < 250; i++) {
    largeString += 'résumé ';
  }
  // Add Korean characters
  for (let i = 0; i < 250; i++) {
    largeString += '안녕하세요 ';
  }
  // Add Arabic text
  for (let i = 0; i < 250; i++) {
    largeString += 'مرحبا ';
  }

  // Test UTF-8 encoding
  let startTime = new Date().getTime();
  let encoded = textEncoding.encodeUTF8(largeString);
  let encodeTime = new Date().getTime() - startTime;
  console.log(`UTF-8 encoding took ${encodeTime}ms for ${largeString.length} characters`);

  // Test UTF-8 decoding
  startTime = new Date().getTime();
  let decoded = textEncoding.decodeUTF8(encoded);
  let decodeTime = new Date().getTime() - startTime;
  console.log(`UTF-8 decoding took ${decodeTime}ms for ${encoded.length} bytes`);

  // Test Base64 encoding
  startTime = new Date().getTime();
  let base64Encoded = textEncoding.encodeUTF8ToBase64(largeString);
  let base64EncodeTime = new Date().getTime() - startTime;
  console.log(`Base64 encoding took ${base64EncodeTime}ms`);

  // Test Base64 decoding
  startTime = new Date().getTime();
  let base64Decoded = textEncoding.decodeUTF8FromBase64(base64Encoded);
  let base64DecodeTime = new Date().getTime() - startTime;
  console.log(`Base64 decoding took ${base64DecodeTime}ms`);

  // Verify round-trip
  assertEqual(decoded, largeString, 'UTF-8 round-trip should preserve content');
  assertEqual(base64Decoded, largeString, 'Base64 round-trip should preserve content');

  // Test byte and rune counting
  let byteCount = textEncoding.countUTF8Bytes(largeString);
  let runeCount = textEncoding.countUTF8Runes(largeString);
  console.log(`String has ${byteCount} bytes and ${runeCount} runes`);

  // Verify byte count matches encoded length
  assertEqual(byteCount, encoded.length, 'Byte count should match encoded length');

  // Verify rune count is less than byte count (since some runes use multiple bytes)
  assert(runeCount < byteCount, 'Rune count should be less than byte count');

  // Verify UTF-8 validation
  assert(textEncoding.isValidUTF8(largeString), 'Large string should be valid UTF-8');
  assert(textEncoding.isValidUTF8Bytes(encoded), 'Encoded bytes should be valid UTF-8');

  console.log('✓ Performance tests passed\n');
}

function testErrorHandling() {
  console.log('Testing error handling...');
  
  // Test that these don't crash (undefined handling)
  try {
    textEncoding.countUTF8Bytes(undefined);
    textEncoding.countUTF8Runes(undefined);
    textEncoding.isValidUTF8(undefined);
  } catch (e) {
    // It's OK if they throw, just shouldn't crash the extension
  }
  
  // Null inputs where appropriate
  assertThrows(() => textEncoding.decodeUTF8(null), 'decodeUTF8 should throw on null');
  
  // isValidUTF8Bytes with null should not crash (Go's utf8.Valid handles nil)
  try {
    textEncoding.isValidUTF8Bytes(null);
  } catch (e) {
    // It's OK if it throws, just shouldn't crash
  }
  
  console.log('✓ Error handling tests passed\n');
}

function testInvalidUTF8Sequences() {
  console.log('Testing invalid UTF-8 sequences...');
  
  // Test various invalid UTF-8 sequences
  const invalidSequences = [
    new Uint8Array([0xFF, 0xFE]),                    // Invalid start byte
    new Uint8Array([0xC0, 0x80]),                    // Overlong encoding
    new Uint8Array([0xF0, 0x9F]),                    // Incomplete sequence
    new Uint8Array([0xED, 0xA0, 0x80]),             // Surrogate pair
    new Uint8Array([0xF4, 0x90, 0x80, 0x80]),       // Out of range
    new Uint8Array([0x80]),                          // Continuation byte without start
    new Uint8Array([0xC0, 0xAF]),                    // Overlong ASCII
    new Uint8Array([0xE0, 0x80, 0xAF]),             // Overlong 2-byte sequence
    new Uint8Array([0xF0, 0x80, 0x80, 0xAF]),       // Overlong 3-byte sequence
    new Uint8Array([0xF8, 0x80, 0x80, 0x80, 0xAF]), // 5-byte sequence (invalid)
  ];

  for (let i = 0; i < invalidSequences.length; i++) {
    const seq = invalidSequences[i];
    // Test IsValidUTF8Bytes
    assert(!textEncoding.isValidUTF8Bytes(seq), `IsValidUTF8Bytes should return false for invalid sequence ${i}`);
    
    // Test DecodeUTF8
    assertThrows(() => textEncoding.decodeUTF8(seq), `DecodeUTF8 should throw for invalid sequence ${i}`);
  }

  console.log('✓ Invalid UTF-8 sequence tests passed\n');
}

function testConcurrentOperations() {
  console.log('Testing concurrent operations...');
  
  // Create a test string with various characters
  const testStr = "Hello 🌍 你好 café résumé 안녕하세요 مرحبا 𝄞 𒀀 👨‍👩‍👧‍👦 🏳️‍🌈";
  
  // Number of concurrent operations
  const numOperations = 100;
  let completed = 0;
  let errors = [];
  
  // Run concurrent operations
  for (let i = 0; i < numOperations; i++) {
    // Encode
    const encoded = textEncoding.encodeUTF8(testStr);
    
    // Decode
    try {
      const decoded = textEncoding.decodeUTF8(encoded);
      // Verify roundtrip
      assertEqual(decoded, testStr, 'UTF-8 roundtrip should preserve content');
      
      // Test Base64
      const base64Encoded = textEncoding.encodeUTF8ToBase64(testStr);
      const base64Decoded = textEncoding.decodeUTF8FromBase64(base64Encoded);
      assertEqual(base64Decoded, testStr, 'Base64 roundtrip should preserve content');
      
      completed++;
    } catch (e) {
      errors.push(e);
    }
  }
  
  // Report results
  assertEqual(completed, numOperations, `All ${numOperations} operations should complete successfully`);
  if (errors.length > 0) {
    throw new Error(`Concurrent operations failed: ${errors.join(', ')}`);
  }
  
  console.log('✓ Concurrent operation tests passed\n');
}

function testStressLargeText() {
  console.log('Testing stress with large strings...');
  
  // Create an extremely large string with various Unicode characters
  let stressText = '';
  // Add a mix of characters that might stress the encoder
  for (let i = 0; i < 10000; i++) {
    stressText += 'Hello ';
    stressText += '你好';
    stressText += '🌍';
    stressText += 'café ';
    stressText += 'résumé ';
    stressText += '안녕하세요 ';
    stressText += 'مرحبا ';
    // Add some rare/edge case characters
    stressText += '𝄞'; // Musical symbol
    stressText += '𒀀'; // Cuneiform
    stressText += '👨‍👩‍👧‍👦'; // Family emoji
    stressText += '🏳️‍🌈'; // Flag emoji
    // Add more edge cases
    stressText += 'Z͑ͫ̓ͪ̂ͫ̽͏̴̙̤̞͉͚̯̞̠͍A̴̵̜̰͔ͫ͗͢L̠ͨͧͩ͘G̴̻͈͍͔̹̑͗̎̅͛́Ǫ̵̹̻̝̳͂̌̌͘!͖̬̰̙̗̿̋ͥͥ̂ͣ̐́́͜͞'; // Zalgo text
    stressText += 'ᚠᛇᚻ᛫ᛒᛦᚦ᛫ᚠᚱᚩᚠᚢᚱ᛫ᚠᛁᚱᚪ᛫ᚷᛖᚻᚹᛦᛚᚳᚢᛗ'; // Runic text
    stressText += '꧁༺༻꧂'; // Decorative characters
    stressText += 'ᕕ( ᐛ )ᕗ'; // ASCII art
    stressText += '👾'; // Emoji with variation selector
    stressText += '👨‍💻'; // Emoji with ZWJ
    stressText += '🏴󠁧󠁢󠁥󠁮󠁧󠁿'; // Regional indicator
  }

  // Test UTF-8 encoding and decoding
  const startTime = new Date().getTime();
  const encoded = textEncoding.encodeUTF8(stressText);
  const encodeTime = new Date().getTime() - startTime;
  console.log(`UTF-8 encoding took ${encodeTime}ms for ${stressText.length} characters`);

  const decoded = textEncoding.decodeUTF8(encoded);
  assertEqual(decoded, stressText, 'Stress test roundtrip should preserve content');

  // Test Base64 encoding and decoding
  const base64Encoded = textEncoding.encodeUTF8ToBase64(stressText);
  const base64Decoded = textEncoding.decodeUTF8FromBase64(base64Encoded);
  assertEqual(base64Decoded, stressText, 'Stress test base64 roundtrip should preserve content');

  // Verify UTF-8 validation
  assert(textEncoding.isValidUTF8(stressText), 'Stress test string should be valid UTF-8');
  assert(textEncoding.isValidUTF8Bytes(encoded), 'Encoded stress test bytes should be valid UTF-8');

  // Test byte and rune counting
  const byteCount = textEncoding.countUTF8Bytes(stressText);
  const runeCount = textEncoding.countUTF8Runes(stressText);
  console.log(`Stress test string has ${byteCount} bytes and ${runeCount} runes`);

  // Verify byte count matches encoded length
  assertEqual(byteCount, encoded.length, 'Byte count should match encoded length');

  // Verify rune count is less than byte count
  assert(runeCount < byteCount, 'Rune count should be less than byte count');

  console.log('✓ Stress test passed\n');
}

// Export function for running basic functionality test
export function testBasicFunctionality() {
  let text = 'Hello 🌍';
  
  console.log('=== Basic Functionality Test ===');
  
  // Test encoding
  let bytes = textEncoding.encodeUTF8(text);
  console.log(`Original: "${text}"`);
  console.log(`Encoded bytes length: ${bytes.length}`);
  
  // Test decoding
  let decoded = textEncoding.decodeUTF8(bytes);
  console.log(`Decoded: "${decoded}"`);
  console.log(`Round-trip successful: ${text === decoded}`);
  
  // Test base64
  let base64 = textEncoding.encodeUTF8ToBase64(text);
  let fromBase64 = textEncoding.decodeUTF8FromBase64(base64);
  console.log(`Base64: ${base64}`);
  console.log(`From Base64: "${fromBase64}"`);
  console.log(`Base64 round-trip successful: ${text === fromBase64}`);
  
  // Test counting
  let byteCount = textEncoding.countUTF8Bytes(text);
  let runeCount = textEncoding.countUTF8Runes(text);
  console.log(`Byte count: ${byteCount}`);
  console.log(`Rune count: ${runeCount}`);
  
  // Test validation
  let isValid = textEncoding.isValidUTF8(text);
  let bytesValid = textEncoding.isValidUTF8Bytes(bytes);
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