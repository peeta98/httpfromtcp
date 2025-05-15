package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	Initialized requestState = iota
	Done
)

const (
	crlf       = "\r\n"
	bufferSize = 8
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	request := &Request{state: Initialized}

	for request.state != Done {
		// Grow buffer if needed
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// Read more data into the buffer starting at readToIndex
		bytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = Done
				break
			}
			return nil, err
		}
		readToIndex += bytesRead

		// Parse the current buffer data
		bytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[bytesParsed:readToIndex])
		readToIndex -= bytesParsed
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	// It needs more data before it can parse the request line.
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	// We add 2 to the total bytes read for the CRLF at the end
	totalBytesRead := len(requestLineText) + 2

	return requestLine, totalBytesRead, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formated request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed HTTP version: %s", parts[2])
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case Initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.state = Done

		return n, nil
	case Done:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}
