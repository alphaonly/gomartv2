package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	HTTPServer *http.Server
}

func NewServer(httpServer *http.Server) *Server {
	return &Server{HTTPServer: httpServer}
}

func (s Server) Run() {
	err := s.HTTPServer.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func (s Server) Stop(ctx context.Context) error {
	time.Sleep(time.Second * 2)
	err := s.HTTPServer.Shutdown(ctx)
	log.Println("Stop http server")
	return err
}
