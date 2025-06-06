package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const serverAddr = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP address %s: %v\n", serverAddr, err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing UDP: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", serverAddr)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading message: %v\n", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", message)
	}
}
