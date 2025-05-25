package response

import (
	"fmt"
	"github.com/peeta98/httpfromtcp/internal/headers"
	"io"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := headers.NewHeaders()

	defaultHeaders.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	defaultHeaders.Set("Connection", "close")
	defaultHeaders.Set("Content-Type", "text/plain")

	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}

	// Write the final CRLF to indicate the end of headers
	_, err := fmt.Fprint(w, "\r\n")
	return err
}
