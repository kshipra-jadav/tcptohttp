package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(fp io.ReadCloser) <-chan string {
	linesChannel := make(chan string)

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
	fp, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	defer fp.Close()

	linesChannel := getLinesChannel(fp)

	for line := range linesChannel {
		fmt.Printf("[LINE]: %v\n", line)
	}

}
