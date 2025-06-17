// Package text_encoding provides text encoding/decoding for xk6 extension
package text_encoding

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"

	"github.com/grafana/sobek"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/text-encoding", New())
}

type (
	TextEncoding struct {
		vu      modules.VU
		exports *sobek.Object
	}
	RootModule struct{}
	Module     struct {
		*TextEncoding
	}
)

var (
	_ modules.Instance = &Module{}
	_ modules.Module   = &RootModule{}
)

// TextEncoder holds the encoding configuration
type TextEncoder struct {
	encoding encoding.Encoding
	label    string
}

// TextDecoder holds the decoding configuration
type TextDecoder struct {
	encoding encoding.Encoding
	label    string
}

// Utils provides utility functions for text encoding
type Utils struct{}

// Buffer pool for reusing memory
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// getBuffer gets a buffer from the pool
func getBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// putBuffer returns a buffer to the pool
func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

// getEncoding returns the appropriate encoding based on the label
func getEncoding(label string) (encoding.Encoding, error) {
	label = strings.ToLower(strings.TrimSpace(label))

	switch label {
	case "utf-8", "utf8":
		return unicode.UTF8, nil
	case "utf-16", "utf16":
		return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM), nil
	case "utf-16le", "utf16le":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), nil
	case "utf-16be", "utf16be":
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM), nil
	case "iso-8859-1", "latin1":
		return charmap.ISO8859_1, nil
	case "iso-8859-2", "latin2":
		return charmap.ISO8859_2, nil
	case "iso-8859-3", "latin3":
		return charmap.ISO8859_3, nil
	case "iso-8859-4", "latin4":
		return charmap.ISO8859_4, nil
	case "iso-8859-5":
		return charmap.ISO8859_5, nil
	case "iso-8859-6":
		return charmap.ISO8859_6, nil
	case "iso-8859-7":
		return charmap.ISO8859_7, nil
	case "iso-8859-8":
		return charmap.ISO8859_8, nil
	case "iso-8859-9", "latin5":
		return charmap.ISO8859_9, nil
	case "iso-8859-10", "latin6":
		return charmap.ISO8859_10, nil
	case "iso-8859-13", "latin7":
		return charmap.ISO8859_13, nil
	case "iso-8859-14", "latin8":
		return charmap.ISO8859_14, nil
	case "iso-8859-15", "latin9":
		return charmap.ISO8859_15, nil
	case "iso-8859-16", "latin10":
		return charmap.ISO8859_16, nil
	case "windows-1250":
		return charmap.Windows1250, nil
	case "windows-1251":
		return charmap.Windows1251, nil
	case "windows-1252":
		return charmap.Windows1252, nil
	case "windows-1253":
		return charmap.Windows1253, nil
	case "windows-1254":
		return charmap.Windows1254, nil
	case "windows-1255":
		return charmap.Windows1255, nil
	case "windows-1256":
		return charmap.Windows1256, nil
	case "windows-1257":
		return charmap.Windows1257, nil
	case "windows-1258":
		return charmap.Windows1258, nil
	case "koi8-r":
		return charmap.KOI8R, nil
	case "koi8-u":
		return charmap.KOI8U, nil
	case "shift-jis", "shift_jis", "sjis":
		return japanese.ShiftJIS, nil
	case "euc-jp", "eucjp":
		return japanese.EUCJP, nil
	case "iso-2022-jp", "iso2022jp":
		return japanese.ISO2022JP, nil
	case "gbk":
		return simplifiedchinese.GBK, nil
	case "gb18030":
		return simplifiedchinese.GB18030, nil
	case "big5":
		return traditionalchinese.Big5, nil
	case "euc-kr", "euckr":
		return korean.EUCKR, nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", label)
	}
}

// New creates a new instance of the root module.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance creates a new instance of the TextEncoding module.
func (*RootModule) NewModuleInstance(virtualUser modules.VU) modules.Instance {
	runtime := virtualUser.Runtime()

	// Create a new TextEncoding module.
	moduleInstance := &Module{
		TextEncoding: &TextEncoding{
			vu:      virtualUser,
			exports: runtime.NewObject(),
		},
	}

	mustExport := func(name string, value interface{}) {
		if err := moduleInstance.exports.Set(name, value); err != nil {
			common.Throw(runtime, err)
		}
	}

	// Export the constructors from the TextEncoding module to the JS code.
	// The TextEncoder is a constructor and must be called with new, e.g. new TextEncoder(...).
	mustExport("TextEncoder", moduleInstance.textEncoderClass)
	// The TextDecoder is a constructor and must be called with new, e.g. new TextDecoder(...).
	mustExport("TextDecoder", moduleInstance.textDecoderClass)
	// The Utils is a constructor and must be called with new, e.g. new Utils().
	mustExport("Utils", moduleInstance.utilsClass)

	return moduleInstance
}

// Exports returns the exports of the TextEncoding module, which are the functions
// that can be called from the JS code.
func (m *Module) Exports() modules.Exports {
	return modules.Exports{
		Default: m.TextEncoding.exports,
	}
}

// textEncoderClass is the JS constructor for TextEncoder
func (m *Module) textEncoderClass(call sobek.ConstructorCall) *sobek.Object {
	rt := m.vu.Runtime()

	var label string
	if len(call.Arguments) > 0 {
		label = call.Arguments[0].String()
	}
	if label == "" {
		label = "utf-8" // Default to UTF-8
	}

	enc, err := getEncoding(label)
	if err != nil {
		common.Throw(rt, err)
	}

	encoder := &TextEncoder{
		encoding: enc,
		label:    label,
	}

	obj := rt.NewObject()
	obj.Set("encode", encoder.Encode)
	obj.Set("encodeString", encoder.EncodeString)
	obj.Set("getEncoding", encoder.GetEncoding)

	return obj
}

// textDecoderClass is the JS constructor for TextDecoder
func (m *Module) textDecoderClass(call sobek.ConstructorCall) *sobek.Object {
	rt := m.vu.Runtime()

	var label string
	if len(call.Arguments) > 0 {
		label = call.Arguments[0].String()
	}
	if label == "" {
		label = "utf-8" // Default to UTF-8
	}

	enc, err := getEncoding(label)
	if err != nil {
		common.Throw(rt, err)
	}

	decoder := &TextDecoder{
		encoding: enc,
		label:    label,
	}

	obj := rt.NewObject()
	obj.Set("decode", decoder.Decode)
	obj.Set("getEncoding", decoder.GetEncoding)

	return obj
}

// utilsClass is the JS constructor for Utils
func (m *Module) utilsClass(call sobek.ConstructorCall) *sobek.Object {
	rt := m.vu.Runtime()

	utils := &Utils{}

	obj := rt.NewObject()
	obj.Set("utf8ByteLength", utils.UTF8ByteLength)
	obj.Set("isValidEncoding", utils.IsValidEncoding)
	obj.Set("getSupportedEncodings", utils.GetSupportedEncodings)

	return obj
}

// Encode encodes a string to bytes using the specified encoding
func (te *TextEncoder) Encode(text string) ([]byte, error) {
	if text == "" {
		return []byte{}, nil
	}

	// For UTF-8, we can optimize by returning the string as bytes directly
	if te.label == "utf-8" || te.label == "utf8" {
		return []byte(text), nil
	}

	// For other encodings, use the encoding package
	encoder := te.encoding.NewEncoder()
	encoded, err := encoder.Bytes([]byte(text))
	if err != nil {
		return nil, fmt.Errorf("failed to encode text: %w", err)
	}

	return encoded, nil
}

// EncodeString is a convenience method that returns the encoded bytes as a string
func (te *TextEncoder) EncodeString(text string) (string, error) {
	encoded, err := te.Encode(text)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

// GetEncoding returns the encoding label
func (te *TextEncoder) GetEncoding() string {
	return te.label
}

// Decode decodes bytes to a string using the specified encoding
func (td *TextDecoder) Decode(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	// For UTF-8, we can optimize by returning the bytes as string directly
	if td.label == "utf-8" || td.label == "utf8" {
		return string(data), nil
	}

	// For other encodings, use the encoding package
	decoder := td.encoding.NewDecoder()
	decoded, err := decoder.Bytes(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode data: %w", err)
	}

	return string(decoded), nil
}

// GetEncoding returns the encoding label
func (td *TextDecoder) GetEncoding() string {
	return td.label
}

// UTF8ByteLength returns the byte length of a string when encoded in UTF-8
// This is much faster than the JavaScript equivalent as it uses Go's optimized UTF-8 handling
func (u *Utils) UTF8ByteLength(str string) int {
	// For UTF-8, the byte length is simply len(str) since Go strings are UTF-8 encoded
	// This is the most efficient way to get UTF-8 byte length in Go
	return len(str)
}

// IsValidEncoding checks if the given encoding label is supported
func (u *Utils) IsValidEncoding(label string) bool {
	_, err := getEncoding(label)
	return err == nil
}

// GetSupportedEncodings returns a list of supported encoding labels
func (u *Utils) GetSupportedEncodings() []string {
	return []string{
		"utf-8", "utf8",
		"utf-16", "utf16", "utf-16le", "utf16le", "utf-16be", "utf16be",
		"iso-8859-1", "latin1",
		"iso-8859-2", "latin2",
		"iso-8859-3", "latin3",
		"iso-8859-4", "latin4",
		"iso-8859-5",
		"iso-8859-6",
		"iso-8859-7",
		"iso-8859-8",
		"iso-8859-9", "latin5",
		"iso-8859-10", "latin6",
		"iso-8859-13", "latin7",
		"iso-8859-14", "latin8",
		"iso-8859-15", "latin9",
		"iso-8859-16", "latin10",
		"windows-1250", "windows-1251", "windows-1252", "windows-1253",
		"windows-1254", "windows-1255", "windows-1256", "windows-1257", "windows-1258",
		"koi8-r", "koi8-u",
		"shift-jis", "shift_jis", "sjis",
		"euc-jp", "eucjp",
		"iso-2022-jp", "iso2022jp",
		"gbk", "gb18030", "big5",
		"euc-kr", "euckr",
	}
}
