package headers

import (
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func parseHeadersLine(h Headers, str string) (err error) {

	firstColonIndex := strings.Index(str, ":")
	if firstColonIndex == -1 {
		return fmt.Errorf("invlaid header")
	}
	key := str[:firstColonIndex]
	value := str[firstColonIndex+1:]
	for _, char := range key {
		if !unicode.IsDigit(char) && !unicode.IsLetter(char) && char != '-' {
			return fmt.Errorf("invalid Header")
		}
	}
	newvalue, exists := h[strings.ToLower(key)]
	if exists {
		h[strings.ToLower(key)] = newvalue + "," + strings.Trim(value, " ")
	} else {
		h[strings.ToLower(key)] = strings.Trim(value, " ")
	}
	return nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	strData := string(data)
	firstHeaderEndIndex := strings.Index(strData, "\r\n")
	if firstHeaderEndIndex == -1 {
		return 0, false, nil
	}
	if firstHeaderEndIndex == 0 {
		return 2, true, nil
	}
	slicedData := strData[:firstHeaderEndIndex]
	consumedData := len(slicedData) + 2
	parseErr := parseHeadersLine(h, slicedData)
	if parseErr != nil {
		return 0, false, parseErr
	}

	return consumedData, false, nil
}
