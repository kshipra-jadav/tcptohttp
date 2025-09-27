package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/kshipra-jadav/tcptohttp/.internal/request"
)


func getLinesChannel(fp io.ReadCloser) <-chan string {
	linesChannel := make(chan string)
	fmt.Println("Connection accepted!")

	go func() {
		buf := make([]byte, 8)
		line := ""
		for {
			n, err := fp.Read(buf)
			if n > 0 {
				str := string(buf[:n])
				splits := strings.Split(str, "\n")
				if len(splits) > 1 {
					line += splits[0]
					linesChannel <- line
					line = strings.Join(splits[1:], "")
				} else {
					line += splits[0]
				}
			}
			if err != nil {
				if line != "" {
					linesChannel <- line
				}
				if err == io.EOF {
					fmt.Println("EOF Signal Received. Closing the TCP Connection.")
					close(linesChannel)
					return
				}
				log.Fatalf("Error: %v", err)
			}
		}
	}()

	return linesChannel
}

func main() {
	req, err := request.RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	if err != nil {
		log.Fatalf("error - %v", err)
	}
	fmt.Print(req)
	return
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Fatalf("Error in accepting: %v", err)
	}
	defer conn.Close()

	if err != nil {
		log.Fatalf("Error in reading: %v", err)
	}

	for line := range getLinesChannel(conn) {
		fmt.Println(line)
	}

	
}
