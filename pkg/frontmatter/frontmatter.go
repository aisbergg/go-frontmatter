/*
Package frontmatter implements detection and decoding for various content
front matter formats.

  The following front matter formats are supported by default.

  - YAML identified by:
    • opening and closing `---` lines.
    • opening `---yaml` and closing `---` lines.
  - TOML identified by:
    • opening and closing `+++` lines.
    • opening `---toml` and closing `---` lines.
  - JSON identified by:
    • opening and closing `;;;` lines.
    • opening `---json` and closing `---` lines.
    • a single JSON object followed by an empty line.

If the default formats are not suitable for your use case, you can easily bring
your own. See the examples for more information.
*/
package frontmatter

import (
	"errors"
	"io"
)

// ErrNotFound is reported by `MustParse` when a front matter is not found.
var ErrNotFound = errors.New("not found")

// Parse decodes the front matter from the specified reader into the value
// pointed to by `v`, and returns the rest of the data. If a front matter
// is not present, the original data is returned and `v` is left unchanged.
// Front matters are detected and decoded based on the passed in `formats`.
// If no formats are provided, the default formats are used.
func Parse(r io.Reader, v interface{}, formats ...*Format) ([]byte, error) {
	return newParser(r).parse(v, formats, false)
}

// ParseWithBuffer is the same as `Parse` but uses a given buffer to read the
// data. This is used to avoid unnecessary allocations.
func ParseWithBuffer(b []byte, r io.Reader, v interface{}, formats ...*Format) ([]byte, error) {
	return newParserWithBuffer(b, r).parse(v, formats, false)
}

// MustParse decodes the front matter from the specified reader into the
// value pointed to by `v`, and returns the rest of the data. If a front
// matter is not present, `ErrNotFound` is reported.
// Front matters are detected and decoded based on the passed in `formats`.
// If no formats are provided, the default formats are used.
func MustParse(r io.Reader, v interface{}, formats ...*Format) ([]byte, error) {
	return newParser(r).parse(v, formats, true)
}

// MustParseWithBuffer is the same as `MustParse` but uses a given buffer to
// read the data. This is used to avoid unnecessary allocations.
func MustParseWithBuffer(b []byte, r io.Reader, v interface{}, formats ...*Format) ([]byte, error) {
	return newParserWithBuffer(b, r).parse(v, formats, true)
}
