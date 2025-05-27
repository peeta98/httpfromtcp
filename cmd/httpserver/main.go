package main

import (
	"errors"
	"fmt"
	"github.com/peeta98/httpfromtcp/internal/request"
	"github.com/peeta98/httpfromtcp/internal/response"
	"github.com/peeta98/httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	reqPath := req.RequestLine.RequestTarget

	if reqPath == "/yourproblem" {
		handler400(w, req)
		return
	}

	if reqPath == "/myproblem" {
		handler500(w, req)
		return
	}

	if strings.HasPrefix(reqPath, "/httpbin") {
		proxyHandler(w, req)
		return
	}

	handler200(w, req)
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(http.StatusBadRequest)
	body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(http.StatusInternalServerError)
	body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(http.StatusOK)
	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	headers := response.GetDefaultHeaders(len(body))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target
	fmt.Println("Proxying to", url)

	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status: %s\n", resp.Status)
		handler500(w, req)
		return
	}

	w.WriteStatusLine(http.StatusOK)

	headers := response.GetDefaultHeaders(0)
	headers.Remove("Content-Length")
	headers.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(headers)

	const maxChunkSize = 1024
	buf := make([]byte, maxChunkSize)

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			fmt.Printf("Read %d bytes: %s\n", n, buf[:n])
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading response body:", err)
			return
		}

		_, err = w.WriteChunkedBody(buf[:n])
		if err != nil {
			fmt.Println("Error writing chunked body:", err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body done:", err)
	}
}
