package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Error establishing udp conn: %v", err)
	}

	udpConn, err := net.DialUDP("udp", nil, conn)
	if err != nil {
		log.Fatalf("Error dailing udp conn: %v", err)
	}
	defer udpConn.Close()

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := r.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading string: %v", err)
		}
		bytesWritten, err := udpConn.Write([]byte(input))
		if err != nil {
			log.Fatalf("Error writing to udp connection: %v", err)
		}
		fmt.Printf("Bytes written: %v\n", bytesWritten)

	}

}
