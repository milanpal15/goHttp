package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(buf []byte) (RequestLine, error) {
	bufToStr := string(buf)
	slicedStr := strings.Split(bufToStr, "\r\n")
	values := strings.Split(slicedStr[0], " ")
	if len(values) != 3 {
		return RequestLine{}, fmt.Errorf("invalid parse request")
	}
	allowed := []string{"GET", "POST"}
	valid := false
	for _, allowedVal := range allowed {
		if allowedVal == values[0] {
			valid = true
		}
	}
	if valid != true {
		return RequestLine{}, fmt.Errorf("invalid parse request")
	}
	if !strings.HasPrefix(values[2], "HTTP/") {
		return RequestLine{}, fmt.Errorf("invalid parse request")
	}
	version := strings.Split(values[2], "/")[1]
	if version != "1.1" {
		return RequestLine{}, fmt.Errorf("invalid parse request")
	}
	return RequestLine{
		Method:        values[0],
		RequestTarget: values[1],
		HttpVersion:   version,
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	readBuf := make([]byte, 8)
	output := []byte{}
	for {
		n, err := reader.Read(readBuf)
		if n > 0 {
			chunk := readBuf[:n]
			output = append(output, chunk...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return &Request{
				RequestLine: RequestLine{},
			}, err
		}
	}
	parsedOutput, err := parseRequestLine(output)
	if err != nil {
		return &Request{
			RequestLine: RequestLine{},
		}, err
	}
	return &Request{
		RequestLine: parsedOutput,
	}, nil

}
