import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  console.log('=== xk6-text-encoding Test ===');
  
  const textEncoding = new TextEncoding();
  
  // Test 1: Basic UTF-8 encoding/decoding
  console.log('\n1. Basic UTF-8 Test:');
  try {
    const utf8Encoder = textEncoding.newTextEncoder("utf-8");
    const utf8Decoder = textEncoding.newTextDecoder("utf-8");
    
    const text = "Hello, 世界! 🌍";
    const encoded = utf8Encoder.encode(text);
    const decoded = utf8Decoder.decode(encoded);
    
    console.log(`Original: ${text}`);
    console.log(`Encoded length: ${encoded.length}`);
    console.log(`Decoded: ${decoded}`);
    console.log(`Match: ${text === decoded}`);
  } catch (error) {
    console.log(`Error in UTF-8 test: ${error.message}`);
  }
  
  // Test 2: Different encodings
  console.log('\n2. Different Encodings Test:');
  
  const encodings = [
    { name: "UTF-8", label: "utf-8" },
    { name: "ISO-8859-1", label: "iso-8859-1" },
    { name: "Windows-1252", label: "windows-1252" },
    { name: "Shift-JIS", label: "shift-jis" }
  ];
  
  const testText = "Hello World!";
  
  for (const enc of encodings) {
    try {
      const encoder = textEncoding.newTextEncoder(enc.label);
      const decoder = textEncoding.newTextDecoder(enc.label);
      
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
    const invalidEncoder = textEncoding.newTextEncoder("invalid-encoding");
    console.log("Should not reach here");
  } catch (error) {
    console.log(`Expected error: ${error.message}`);
  }
  
  // Test 4: Empty string handling
  console.log('\n4. Empty String Test:');
  try {
    const emptyEncoder = textEncoding.newTextEncoder("utf-8");
    const emptyDecoder = textEncoding.newTextDecoder("utf-8");
    
    const emptyEncoded = emptyEncoder.encode("");
    const emptyDecoded = emptyDecoder.decode(emptyEncoded);
    
    console.log(`Empty string encoded length: ${emptyEncoded.length}`);
    console.log(`Empty string decoded: "${emptyDecoded}"`);
  } catch (error) {
    console.log(`Error in empty string test: ${error.message}`);
  }
  
  // Test 5: Special characters
  console.log('\n5. Special Characters Test:');
  try {
    const utf8Encoder = textEncoding.newTextEncoder("utf-8");
    const utf8Decoder = textEncoding.newTextDecoder("utf-8");
    
    const specialText = "Special chars: áéíóú ñ ç ß € ¥ £";
    const specialEncoded = utf8Encoder.encode(specialText);
    const specialDecoded = utf8Decoder.decode(specialEncoded);
    
    console.log(`Special text: ${specialText}`);
    console.log(`Special encoded length: ${specialEncoded.length}`);
    console.log(`Special decoded: ${specialDecoded}`);
  } catch (error) {
    console.log(`Error in special characters test: ${error.message}`);
  }
  
  // Test 6: Binary data
  console.log('\n6. Binary Data Test:');
  try {
    const utf8Decoder = textEncoding.newTextDecoder("utf-8");
    const binaryData = new Uint8Array([72, 101, 108, 108, 111]); // "Hello"
    const binaryDecoded = utf8Decoder.decode(binaryData);
    
    console.log(`Binary data: ${binaryData}`);
    console.log(`Binary decoded: "${binaryDecoded}"`);
  } catch (error) {
    console.log(`Error in binary data test: ${error.message}`);
  }
  
  // Test 7: UTF-8 Byte Length
  console.log('\n7. UTF-8 Byte Length Test:');
  try {
    const testStrings = ["Hello", "Hello, 世界!", "🌍", "áéíóú"];
    
    for (const str of testStrings) {
      const byteLength = textEncoding.utf8ByteLength(str);
      console.log(`"${str}": ${byteLength} bytes`);
    }
  } catch (error) {
    console.log(`Error in UTF-8 byte length test: ${error.message}`);
  }
  
  // Test 8: Utility functions
  console.log('\n8. Utility Functions Test:');
  try {
    console.log(`UTF-8 supported: ${textEncoding.isValidEncoding("utf-8")}`);
    console.log(`Invalid encoding supported: ${textEncoding.isValidEncoding("invalid")}`);
    
    const supported = textEncoding.getSupportedEncodings();
    console.log(`Supported encodings count: ${supported.length}`);
  } catch (error) {
    console.log(`Error in utility functions test: ${error.message}`);
  }
  
  console.log('\n=== Test Complete ===');
} 