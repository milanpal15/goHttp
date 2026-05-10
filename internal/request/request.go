package request

import (
	"fmt"
	"io"
	"myhttp/internal/headers"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       string
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(buf []byte) (n int, req RequestLine, err error) {
	bufToStr := string(buf)
	slicedStr := strings.Split(bufToStr, "\r\n")
	consumedBytes := len(slicedStr[0]) + 2
	values := strings.Split(slicedStr[0], " ")
	if len(values) != 3 {
		return 0, RequestLine{}, fmt.Errorf("invalid parse request")
	}
	allowed := []string{"GET", "POST"}
	valid := false
	for _, allowedVal := range allowed {
		if allowedVal == values[0] {
			valid = true
		}
	}
	if valid != true {
		return 0, RequestLine{}, fmt.Errorf("invalid parse request")
	}
	if !strings.HasPrefix(values[2], "HTTP/") {
		return 0, RequestLine{}, fmt.Errorf("invalid parse request")
	}
	version := strings.Split(values[2], "/")[1]
	if version != "1.1" {
		return 0, RequestLine{}, fmt.Errorf("invalid parse request")
	}
	return consumedBytes, RequestLine{
		Method:        values[0],
		RequestTarget: values[1],
		HttpVersion:   version,
	}, nil
}

func (r *Request) parse(data []byte) error {
	totalBytesParsed := 0
	for r.state != "done" {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.state == "initialized" {
		n, parsedData, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		r.RequestLine = parsedData
		r.state = "parsingHeaders"
		return n, nil
	}
	if r.state == "parsingHeaders" {
		n, isCompleted, err := r.Headers.Parse(data)
		if isCompleted {
			r.state = "done"
		}
		if err != nil {
			return 0, err
		}
		return n, nil
	}
	return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	RequestObj := Request{
		RequestLine: RequestLine{},
		Headers:     headers.NewHeaders(),
		state:       "initialized",
	}
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
	err := RequestObj.parse(output)
	if err != nil {
		return &RequestObj, err
	}
	return &RequestObj, nil
}
