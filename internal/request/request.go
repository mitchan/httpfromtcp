package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	Method        string
	RequestTarget string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rl, err := parseRequestLine(r)
	if err != nil {
		return nil, err
	}

	request := &Request{
		RequestLine: *rl,
	}

	return request, nil
}

func parseRequestLine(bytes []byte) (*RequestLine, error) {
	content := string(bytes)

	lines := strings.Split(content, "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("no content")
	}

	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid request line")
	}

	if parts[2] != "HTTP/1.1" {
		return nil, errors.New("invalid http version: " + parts[2])
	}

	methodPattern := "^[A-Z]+$"
	regex, err := regexp.Compile(methodPattern)
	if err != nil {
		return nil, errors.New("cannot compile regex")
	}
	if !regex.MatchString(parts[0]) {
		return nil, errors.New("invalid method" + parts[0])
	}

	rl := &RequestLine{
		HttpVersion:   "1.1",
		Method:        parts[0],
		RequestTarget: parts[1],
	}

	return rl, nil
}
