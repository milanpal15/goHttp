package response

import (
	"fmt"
	"io"
)

func WriteResp(w io.Writer) {
	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!"
	_, err := w.Write([]byte(response))
	if err != nil {
		fmt.Print(err)
	}
}
