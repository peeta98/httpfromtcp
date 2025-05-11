package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
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

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println(line)
		}

		fmt.Printf("Connection to %s closed\n", conn.RemoteAddr())
	}
}

func getLinesChannel(conn net.Conn) <-chan string {
	lines := make(chan string)

	go func() {
		defer conn.Close()
		defer close(lines)

		currentLine := ""

		for {
			buf := make([]byte, 8)
			n, err := conn.Read(buf)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
				}

				if errors.Is(err, io.EOF) {
					break
				}

				fmt.Printf("error: %s\n", err.Error())
				break
			}

			str := string(buf[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- currentLine + parts[i]
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return lines
}
