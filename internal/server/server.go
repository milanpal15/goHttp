package server

import (
	"fmt"
	"myhttp/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func HandleError(err error) {

}

func Serve(port int) (*Server, error) {
	newServer := Server{}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	newServer.listener = listen
	go newServer.listen()
	return &newServer, nil
}

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
	defer conn.Close()
	response.WriteResp(conn)
}
