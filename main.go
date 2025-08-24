package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	currentLine := ""

	for {
		bytes := make([]byte, 8)
		n, err := file.Read(bytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				os.Exit(0)
			}

			log.Fatal(err)
		}

		read := string(bytes[:n])
		lines := strings.Split(read, "\n")

		if len(lines) == 1 {
			currentLine += lines[0]
			continue
		}

		for i, line := range lines {
			if i != len(lines)-1 {
				currentLine += line
				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
			} else {
				currentLine = line
			}
		}
	}

	fmt.Printf("read: %s\n", currentLine)
}
