package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		currentLine := ""

		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
					break
				}

				if errors.Is(err, io.EOF) {
					break
				}

				log.Fatal(err)
			}

			read := string(buffer[:n])
			lines := strings.Split(read, "\n")

			for i := 0; i < len(lines)-1; i++ {
				ch <- fmt.Sprintf("%s%s", currentLine, lines[i])
				currentLine = ""
			}
			currentLine += lines[len(lines)-1]
		}
	}()

	return ch
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}
