package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type parserStatus string

const (
	StatusInitialized parserStatus = "init"
	StatusDone parserStatus = "done"
)

type Request struct {
	RequestLine RequestLine
	Status parserStatus
}

type RequestLine struct {
	Method        string
	HttpVersion   string
	RequestTarget string
}

var ErrMissingHeader = errors.New("less than 3 parts in the request header")
var ErrMalformedHTTPMethod = errors.New("invalid HTTP method")
var ErrIncorrectHTTPVersion = errors.New("invalid HTTP version. Only HTTP/1.1 is supported")


var VALID_VERBS = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
const BUFFER_SIZE = 1024


func (r *Request) parse(data []byte) (int, error) {
	read := 0
	outer:
		for {
			switch r.Status {
			case StatusInitialized:
				requestLine, n, err := parseRequestLine(data[read: ])
				if err != nil {
					return 0, err
				}
				if n == 0 {
					break outer
				}
				r.RequestLine = *requestLine
				read += n
				r.Status = StatusDone
			case StatusDone:
				break outer
			}			
		}
	return read, nil
}

func parseRequestLine(req []byte) (*RequestLine, int, error) {
	line := string(req)
	if !strings.Contains(line, "\r\n") {
		return nil, 0, nil
	}
	line = strings.Split(line, "\r\n")[0]

	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return nil, 0, ErrMissingHeader
	}

	httpMethod := parts[0]
	requestTarget := parts[1]
	httpVersion := strings.Split(parts[2], "/")[1]

	fmt.Printf("method: %v\ntarget: %v\nversion: %v\n", httpMethod, requestTarget, httpVersion)

	if !slices.Contains(VALID_VERBS, httpMethod) {
		return nil, 0, ErrMalformedHTTPMethod
	}
	if httpVersion != "1.1" {
		return nil, 0, ErrIncorrectHTTPVersion
	}

	return &RequestLine {
		Method: httpMethod,
		HttpVersion: httpVersion,
		RequestTarget: requestTarget,
	}, len(req), nil
}


func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{Status: StatusInitialized}
	buf := make([]byte, BUFFER_SIZE)
	bytesRead := 0

	for request.Status != StatusDone {
		numBytes, err := reader.Read(buf[bytesRead: ])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		bytesRead += numBytes

		bytesParsed, err := request.parse(buf[: bytesRead])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[bytesParsed: bytesRead])
		bytesRead -= bytesParsed
	}

	fmt.Println(request)

	return request, nil
}