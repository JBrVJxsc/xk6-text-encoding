import { check } from 'k6';
import { TextEncoding } from 'k6/x/text-encoding';

export default function() {
    const textEncoding = new TextEncoding();
    
    const testCases = [
        { input: '', expected: 0, name: 'empty string' },
        { input: 'a', expected: 1, name: 'ASCII character' },
        { input: '©', expected: 2, name: 'Latin-1 Supplement character' },
        { input: 'ह', expected: 3, name: 'Devanagari character' },
        { input: '😊', expected: 4, name: 'Emoji (surrogate pair)' },
        { input: '가', expected: 3, name: 'Korean character (Hangul)' },
        { input: '汉', expected: 3, name: 'Chinese character (CJK)' },
        { input: 'Hello©ह😊', expected: 5 + 2 + 3 + 4, name: 'Mixed characters' },
        { input: '안녕하세요你好', expected: (5 * 3) + (2 * 3), name: 'Mixed Korean and Chinese' },
    ];

    for (const { input, expected, name } of testCases) {
        const len = textEncoding.utf8ByteLength(input);
        check(len, {
            [`${name}`]: () => len === expected,
        });
    }
} 