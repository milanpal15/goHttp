package main

import (
	"log"
	"myhttp/internal/request"
	"myhttp/internal/response"
	"myhttp/internal/server"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func handler(w *response.Response, req *request.Request) {
	target := req.RequestLine.RequestTarget
	w.Headers.Set("Content-Type", "text/html")
	w.HandlerRes.StatusCode = 200
	w.HandlerRes.Message = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	if target == "/yourproblem" {
		w.HandlerRes.StatusCode = 400
		w.HandlerRes.Message = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	}
	if target == "/myproblem" {
		w.HandlerRes.StatusCode = 500
		w.HandlerRes.Message = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	}
}

func main() {
	server, err := server.Serve(handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
