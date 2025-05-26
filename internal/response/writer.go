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
