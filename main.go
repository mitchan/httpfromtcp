package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection accepted")

		for line := range getLinesChannel(connection) {
			fmt.Printf("read: %s\n", line)
		}
	}
}
