package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}
	defer f.Close()

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	currentLine := ""

	for {
		buf := make([]byte, 8)
		n, err := f.Read(buf)
		if err != nil {
			if currentLine != "" {
				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
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
			fmt.Printf("read: %s\n", currentLine+parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}
}
