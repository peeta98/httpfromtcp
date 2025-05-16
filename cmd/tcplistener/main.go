package main

import (
	"fmt"
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

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)

		fmt.Printf("Connection to %s closed\n", conn.RemoteAddr())
	}
}
