import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  console.log('=== Simple Test ===');
  
  try {
    console.log('Creating TextEncoding instance...');
    const textEncoding = new TextEncoding();
    console.log('TextEncoding created successfully!');
    
    console.log('Testing utf8ByteLength...');
    const result = textEncoding.utf8ByteLength("Hello");
    console.log(`Result: ${result}`);
    
    console.log('Testing getSupportedEncodings...');
    const supported = textEncoding.getSupportedEncodings();
    console.log(`Supported encodings: ${supported.length} found`);
    
    console.log('=== Test Passed ===');
  } catch (error) {
    console.log('Error:', error.message);
    console.log('Stack:', error.stack);
  }
} 