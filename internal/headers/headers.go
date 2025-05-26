package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
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
	key := strings.ToLower(string(parts[0]))

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	key = strings.TrimSpace(key)

	if !h.isValidTokens(key) {
		return 0, false, fmt.Errorf("invalid header token found: %s", key)
	}

	if len(parts) != 2 {
		return 0, false, fmt.Errorf("malformed header line: %s", string(rawHeaderLine))
	}

	value := strings.TrimSpace(string(parts[1]))

	h.Set(key, value)
	return idx + 2, false, nil
}

func (h Headers) isValidTokens(value string) bool {
	for _, r := range strings.TrimSpace(value) {
		if !isTChar(r) {
			return false
		}
	}
	return true
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		h[key] = fmt.Sprintf("%s, %s", v, value)
		return
	}

	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	val, ok := h[key]
	return val, ok
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func isTChar(r rune) bool {
	switch {
	case unicode.IsLetter(r), unicode.IsDigit(r):
		return true
	case strings.ContainsRune("!#$%&'*+-.^_`|~", r):
		return true
	default:
		return false
	}
}
