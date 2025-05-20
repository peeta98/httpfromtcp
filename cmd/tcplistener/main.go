package main

import (
	"fmt"
	"github.com/peeta98/httpfromtcp/internal/headers"
	"github.com/peeta98/httpfromtcp/internal/request"
	"log"
	"net"
	"os"
)

const port = ":42069"

func main() {
	tcpListener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer tcpListener.Close()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			// Using log.Fatal() makes deferred functions not run!
			fmt.Printf("\"could not wait for a connection: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("A connection has been accepted from %s\n", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error parsing incoming request: %s\n", err)
			os.Exit(1)
		}

		printRequestLine(req.RequestLine)
		printHeaders(req.Headers)

		fmt.Printf("Connection to %s closed\n", conn.RemoteAddr())
	}
}

func printRequestLine(rl request.RequestLine) {
	fmt.Println("Request line:")
	fmt.Println("- Method:", rl.Method)
	fmt.Println("- Target:", rl.RequestTarget)
	fmt.Println("- Version:", rl.HttpVersion)
}

func printHeaders(h headers.Headers) {
	fmt.Println("Headers: ")
	for k, v := range h {
		fmt.Printf("- %s: %s\n", k, v)
	}
}
