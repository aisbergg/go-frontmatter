package frontmatter

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

type parser struct {
	reader *bufio.Reader
	output *bytes.Buffer

	read  int
	start int
	end   int
}

func newParser(r io.Reader) *parser {
	if _, ok := r.(*bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}
	br := r.(*bufio.Reader)
	return &parser{
		reader: br,
		output: bytes.NewBuffer(make([]byte, 0, 4096)),
	}
}

func newParserWithBuffer(b []byte, r io.Reader) *parser {
	if _, ok := r.(*bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}
	br := r.(*bufio.Reader)
	return &parser{
		reader: br,
		output: bytes.NewBuffer(b[:0]),
	}
}

func (p *parser) parse(v interface{}, formats []*Format, mustParse bool) ([]byte, error) {
	// If no formats are provided, use the default ones.
	if len(formats) == 0 {
		formats = defaultFormats()
	}

	// Detect format.
	f, err := p.detectStart(formats)
	if err != nil {
		return nil, err
	}

	// Extract front matter.
	found := f != nil
	if found {
		if found, err = p.extract(f, v); err != nil {
			return nil, err
		}
	}
	if mustParse && !found {
		return nil, ErrNotFound
	}

	// Read remaining data.
	if _, err := p.output.ReadFrom(p.reader); err != nil {
		return nil, err
	}

	return p.output.Bytes()[p.end:], nil
}

func (p *parser) detectStart(formats []*Format) (*Format, error) {
	for {
		read := p.read

		line, atEOF, err := p.readLine()
		if err != nil || atEOF {
			return nil, err
		}
		if len(line) == 0 {
			continue
		}

		for _, f := range formats {
			if f.Start == string(line) {
				if !f.UnmarshalDelims {
					read = p.read
				}

				p.start = read
				return f, nil
			}
		}

		return nil, nil
	}
}

func (p *parser) extract(f *Format, v interface{}) (bool, error) {
	for {
		read := p.read

		line, atEOF, err := p.readLine()
		if err != nil {
			return false, err
		}

	CheckLine:
		if string(line) != f.End {
			if atEOF {
				return false, err
			}
			continue
		}
		if f.RequiresNewLine {
			if line, atEOF, err = p.readLine(); err != nil {
				return false, err
			}
			if len(line) > 0 {
				goto CheckLine
			}
		}
		if f.UnmarshalDelims {
			read = p.read
		}

		if err := f.Unmarshal(p.output.Bytes()[p.start:read], v); err != nil {
			return false, err
		}

		p.end = p.read
		return true, nil
	}
}

func (p *parser) readLine() ([]byte, bool, error) {
	line, err := p.reader.ReadBytes('\n')

	atEOF := err == io.EOF
	if err != nil && !atEOF {
		return nil, false, err
	}

	p.read += len(line)
	_, err = p.output.Write(line)
	return trimRightSpace(line), atEOF, err
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

// trimRightSpace trims trailing whitespace from the given byte slice.
func trimRightSpace(s []byte) []byte {
	stop := len(s)
	for ; stop > 0; stop-- {
		c := s[stop-1]
		if c >= utf8.RuneSelf {
			return bytes.TrimRightFunc(s[0:stop], unicode.IsSpace)
		}
		// Fast path for ASCII: look for the last ASCII non-space byte
		if asciiSpace[c] == 0 {
			break
		}
	}
	return s[0:stop]
}
