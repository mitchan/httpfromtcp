package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {
		bytes := make([]byte, 8)
		_, err := file.Read(bytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				os.Exit(0)
			}

			log.Fatal(err)
		}

		fmt.Printf("read: %s\n", bytes)
	}
}
