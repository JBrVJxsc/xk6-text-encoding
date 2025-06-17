import { TextEncoding } from 'k6/x/text-encoding';

export default function () {
  console.log('=== Debug Test ===');
  
  try {
    const textEncoding = new TextEncoding();
    console.log('TextEncoding object created successfully');
    console.log('Available methods:', Object.getOwnPropertyNames(Object.getPrototypeOf(textEncoding)));
    
    // Try to call a simple method
    const supported = textEncoding.getSupportedEncodings();
    console.log('Supported encodings:', supported);
    
  } catch (error) {
    console.log('Error:', error.message);
    console.log('Stack:', error.stack);
  }
} 