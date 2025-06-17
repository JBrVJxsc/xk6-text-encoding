import { TextEncoder, TextDecoder, Utils } from 'k6/x/text-encoding';

export default function () {
  console.log('=== xk6-text-encoding Test ===');
  
  // Test 1: Basic UTF-8 encoding/decoding
  console.log('\n1. Basic UTF-8 Test:');
  try {
    const encoder = new TextEncoder("utf-8");
    const decoder = new TextDecoder("utf-8");
    
    const text = "Hello, 世界! 🌍";
    const encoded = encoder.encode(text);
    const decoded = decoder.decode(encoded);
    
    console.log(`Original: ${text}`);
    console.log(`Encoded length: ${encoded.length}`);
    console.log(`Decoded: ${decoded}`);
    console.log(`Match: ${text === decoded}`);
    console.log(`Encoder encoding: ${encoder.getEncoding()}`);
    console.log(`Decoder encoding: ${decoder.getEncoding()}`);
  } catch (error) {
    console.log(`Error in UTF-8 test: ${error.message}`);
  }
  
  // Test 2: Different encodings
  console.log('\n2. Different Encodings Test:');
  
  const encodings = [
    { name: "UTF-8", label: "utf-8" },
    { name: "UTF-16", label: "utf-16" },
    { name: "ISO-8859-1", label: "iso-8859-1" },
    { name: "Windows-1252", label: "windows-1252" },
    { name: "Shift-JIS", label: "shift-jis" }
  ];
  
  const testText = "Hello World!";
  
  for (const enc of encodings) {
    try {
      const encoder = new TextEncoder(enc.label);
      const decoder = new TextDecoder(enc.label);
      
      const encoded = encoder.encode(testText);
      const decoded = decoder.decode(encoded);
      
      console.log(`${enc.name}: ${encoded.length} bytes, decoded: "${decoded}"`);
    } catch (error) {
      console.log(`${enc.name}: Error - ${error.message}`);
    }
  }
  
  // Test 3: Error handling
  console.log('\n3. Error Handling Test:');
  
  try {
    const invalidEncoder = new TextEncoder("invalid-encoding");
    console.log("Should not reach here");
  } catch (error) {
    console.log(`Expected error for invalid encoding: ${error.message}`);
  }
  
  // Test 4: Empty string handling
  console.log('\n4. Empty String Test:');
  try {
    const encoder = new TextEncoder("utf-8");
    const decoder = new TextDecoder("utf-8");
    
    const emptyEncoded = encoder.encode("");
    const emptyDecoded = decoder.decode(emptyEncoded);
    
    console.log(`Empty string encoded length: ${emptyEncoded.length}`);
    console.log(`Empty string decoded: "${emptyDecoded}"`);
  } catch (error) {
    console.log(`Error in empty string test: ${error.message}`);
  }
  
  // Test 5: Special characters
  console.log('\n5. Special Characters Test:');
  try {
    const encoder = new TextEncoder("utf-8");
    const decoder = new TextDecoder("utf-8");
    
    const specialText = "Special chars: áéíóú ñ ç ß € ¥ £";
    const specialEncoded = encoder.encode(specialText);
    const specialDecoded = decoder.decode(specialEncoded);
    
    console.log(`Special text: ${specialText}`);
    console.log(`Special encoded length: ${specialEncoded.length}`);
    console.log(`Special decoded: ${specialDecoded}`);
  } catch (error) {
    console.log(`Error in special characters test: ${error.message}`);
  }
  
  // Test 6: Binary data
  console.log('\n6. Binary Data Test:');
  try {
    const decoder = new TextDecoder("utf-8");
    const binaryData = new Uint8Array([72, 101, 108, 108, 111]); // "Hello"
    const binaryDecoded = decoder.decode(binaryData);
    
    console.log(`Binary data: [${Array.from(binaryData).join(', ')}]`);
    console.log(`Binary decoded: "${binaryDecoded}"`);
  } catch (error) {
    console.log(`Error in binary data test: ${error.message}`);
  }
  
  // Test 7: UTF-8 Byte Length (Utils class)
  console.log('\n7. UTF-8 Byte Length Test (Utils):');
  try {
    const utils = new Utils();
    const testStrings = [
      "Hello",           // 5 bytes
      "Hello, 世界!",     // 13 bytes (7 ASCII + 6 for Chinese chars)
      "🌍",              // 4 bytes (emoji)
      "áéíóú",           // 10 bytes (5 accented chars, 2 bytes each)
      "",                // 0 bytes
      "A",               // 1 byte
      "AB",              // 2 bytes
      "ABC"              // 3 bytes
    ];
    
    for (const str of testStrings) {
      const byteLength = utils.utf8ByteLength(str);
      console.log(`"${str}": ${byteLength} bytes`);
    }
  } catch (error) {
    console.log(`Error in UTF-8 byte length test: ${error.message}`);
  }
  
  // Test 8: Utility functions
  console.log('\n8. Utility Functions Test:');
  try {
    const utils = new Utils();
    
    // Test isValidEncoding
    const validEncodings = ["utf-8", "UTF-8", "iso-8859-1", "windows-1252", "shift-jis"];
    const invalidEncodings = ["invalid", "", "not-real", "fake-encoding"];
    
    console.log("Valid encodings:");
    for (const encoding of validEncodings) {
      const isValid = utils.isValidEncoding(encoding);
      console.log(`  ${encoding}: ${isValid}`);
    }
    
    console.log("Invalid encodings:");
    for (const encoding of invalidEncodings) {
      const isValid = utils.isValidEncoding(encoding);
      console.log(`  "${encoding}": ${isValid}`);
    }
    
    // Test getSupportedEncodings
    const supported = utils.getSupportedEncodings();
    console.log(`Supported encodings count: ${supported.length}`);
    console.log(`First few supported: ${supported.slice(0, 5).join(', ')}`);
  } catch (error) {
    console.log(`Error in utility functions test: ${error.message}`);
  }
  
  // Test 9: EncodeString method
  console.log('\n9. EncodeString Method Test:');
  try {
    const encoder = new TextEncoder("utf-8");
    
    const testTexts = ["Hello", "世界", "🌍", ""];
    
    for (const text of testTexts) {
      const encodedString = encoder.encodeString(text);
      console.log(`"${text}" encoded as string: "${encodedString}"`);
    }
  } catch (error) {
    console.log(`Error in encodeString test: ${error.message}`);
  }
  
  // Test 10: Default encoding (no parameter)
  console.log('\n10. Default Encoding Test:');
  try {
    const defaultEncoder = new TextEncoder(); // Should default to UTF-8
    const defaultDecoder = new TextDecoder(); // Should default to UTF-8
    
    const text = "Default test 🌍";
    const encoded = defaultEncoder.encode(text);
    const decoded = defaultDecoder.decode(encoded);
    
    console.log(`Default encoder encoding: ${defaultEncoder.getEncoding()}`);
    console.log(`Default decoder encoding: ${defaultDecoder.getEncoding()}`);
    console.log(`Text: ${text}`);
    console.log(`Roundtrip successful: ${text === decoded}`);
  } catch (error) {
    console.log(`Error in default encoding test: ${error.message}`);
  }
  
  // Test 11: Multiple encodings roundtrip
  console.log('\n11. Multiple Encodings Roundtrip Test:');
  
  const roundtripEncodings = [
    "utf-8",
    "utf-16",
    "iso-8859-1",
    "windows-1252"
  ];
  
  for (const encoding of roundtripEncodings) {
    try {
      const encoder = new TextEncoder(encoding);
      const decoder = new TextDecoder(encoding);
      
      // Use simple ASCII text for non-UTF encodings
      const text = encoding.startsWith("utf") ? "Hello 世界 🌍" : "Hello World!";
      
      const encoded = encoder.encode(text);
      const decoded = decoder.decode(encoded);
      
      const success = text === decoded;
      console.log(`${encoding}: ${success ? 'SUCCESS' : 'FAILED'} (${encoded.length} bytes)`);
      
      if (!success) {
        console.log(`  Original: "${text}"`);
        console.log(`  Decoded:  "${decoded}"`);
      }
    } catch (error) {
      console.log(`${encoding}: ERROR - ${error.message}`);
    }
  }
  
  // Test 12: Large text performance
  console.log('\n12. Large Text Performance Test:');
  try {
    const utils = new Utils();
    const encoder = new TextEncoder("utf-8");
    const decoder = new TextDecoder("utf-8");
    
    // Create a large string
    const baseText = "This is a test string with some unicode: 世界 🌍 ";
    const largeText = baseText.repeat(1000); // About 50KB
    
    const startTime = Date.now();
    
    const byteLength = utils.utf8ByteLength(largeText);
    const encoded = encoder.encode(largeText);
    const decoded = decoder.decode(encoded);
    
    const endTime = Date.now();
    
    console.log(`Large text length: ${largeText.length} chars`);
    console.log(`UTF-8 byte length: ${byteLength}`);
    console.log(`Encoded length: ${encoded.length}`);
    console.log(`Roundtrip successful: ${largeText === decoded}`);
    console.log(`Time taken: ${endTime - startTime}ms`);
  } catch (error) {
    console.log(`Error in large text test: ${error.message}`);
  }
  
  console.log('\n=== Test Complete ===');
}