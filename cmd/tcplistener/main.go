package main

import (
	"fmt"
	"io"
	"log"
	"net"
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
	listen, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Failed listener: %s", err)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		out := readchannel(conn)
		for line := range out {
			fmt.Println(line)
		}
	}
}
