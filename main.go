package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func worker(file io.ReadCloser, ch chan string) {
	defer close(ch)
	buffer := make([]byte, 8)
	var buf []byte
	for {
		n, err := file.Read(buffer)

		if n > 0 {
			for _, b := range buffer[:n] {
				if b == '\n' {
					ch <- string(buf)
					buf = buf[:0]
				} else {
					buf = append(buf, b)
				}
			}
		}

		if err == io.EOF {
			// flush remaining buffer if any
			if len(buf) > 0 {
				ch <- string(buf)
			}
			return
		}

		if err != nil {
			ch <- err.Error()
			return
		}
	}
}

func readchannel(file io.ReadCloser) <-chan string {

	ch := make(chan string)
	// return ch
	go worker(file, ch)
	return ch
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()
	out := readchannel(file)

	for line := range out {
		fmt.Println("read:", line)
	}
}
