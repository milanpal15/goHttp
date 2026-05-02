package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	buffer := make([]byte, 8)
	var buf []byte

	for {
		n, err := file.Read(buffer)

		if n > 0 {
			for _, b := range buffer[:n] {
				if b == '\n' {
					fmt.Printf("read: %s\n", string(buf))
					buf = buf[:0]
				} else {
					buf = append(buf, b)
				}
			}
		}

		if err == io.EOF {
			// flush remaining buffer if any
			if len(buf) > 0 {
				fmt.Printf("read: %s\n", string(buf))
			}
			break
		}

		if err != nil {
			fmt.Println("Error reading File:", err)
			return
		}
	}
}

