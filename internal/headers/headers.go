package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const CRLF = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(CRLF))
	// If idx is -1, this means that we still have more data to parse.
	if idx == -1 {
		return 0, false, nil
	}
	// If idx is 0, this means that we have officially finished parsing the headers.
	if idx == 0 {
		return 2, true, nil
	}

	rawHeaderLine := data[:idx]
	parts := bytes.SplitN(rawHeaderLine, []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)
	h[key] = string(value)

	return idx + 2, false, nil
}
