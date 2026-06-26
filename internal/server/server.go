package server

import (
	"fmt"
	"myhttp/internal/headers"
	"myhttp/internal/request"
	"myhttp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

const port = 42069

type Handler func(w *response.Response, req *request.Request)

func Serve(handler Handler) (*Server, error) {
	newServer := Server{}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	newServer.listener = listen
	newServer.handler = handler
	go newServer.listen()
	return &newServer, nil

}

// func Serve(port int) (*Server, error) {
// 	newSer ver := Server{}
// 	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
// 	if err != nil {
// 		return nil, err
// 	}
// 	newServer.listener = listen
// 	go newServer.listen()
// 	return &newServer, nil
// }

func (s *Server) Close() error {
	if !s.closed.Load() {
		s.closed.Store(true)
		err := s.listener.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) listen() {
	for {
		if s.closed.Load() {
			return
		}
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			} else {
				fmt.Printf("%s", err)
			}
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	handleRes := response.Response{}
	handleRes.Writer = conn
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("%s", err)
	}
	handleRes.Headers = headers.NewHeaders()
	defer conn.Close()
	s.handler(&handleRes, req)
	response.HandleRes(&handleRes, handleRes.HandlerRes)
}
