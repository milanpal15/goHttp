package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	remoteAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Failed listener: %s", err)
	}
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error is here: %s", err)
			continue
		}
		conn.Write([]byte(line))
	}
}
