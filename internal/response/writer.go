package response

import (
	"fmt"
	"github.com/peeta98/httpfromtcp/internal/headers"
	"io"
)

type writerState int

const (
	StatusLineState writerState = iota
	HeadersState
	BodyState
	TrailersState
)

type Writer struct {
	writer      io.Writer
	writerState writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer:      w,
		writerState: StatusLineState,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != StatusLineState {
		return fmt.Errorf("cannot write status line: expected state StatusLineState, but current state is %v", w.writerState)
	}

	err := WriteStatusLine(w.writer, statusCode)
	if err != nil {
		return err
	}
	w.writerState = HeadersState

	return nil
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != HeadersState {
		return fmt.Errorf("cannot write headers: expected state HeadersState, but current state is %v", w.writerState)
	}

	err := WriteHeaders(w.writer, h)
	if err != nil {
		return err
	}
	w.writerState = BodyState

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != BodyState {
		return 0, fmt.Errorf("cannot write body: expected state BodyState, but current state is %v", w.writerState)
	}

	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != BodyState {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	var nTotal int
	chunkSizeHex := []byte(fmt.Sprintf("%x\r\n", len(p)))
	n, err := w.writer.Write(chunkSizeHex)
	if err != nil {
		return 0, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	return nTotal + n, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != BodyState {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	finalChunkPart := []byte("0\r\n")
	n, err := w.writer.Write(finalChunkPart)
	if err != nil {
		return n, err
	}
	w.writerState = TrailersState

	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != TrailersState {
		return fmt.Errorf("cannot write trailers in state %d", w.writerState)
	}

	for k, v := range h {
		headerResponse := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.writer.Write([]byte(headerResponse))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
