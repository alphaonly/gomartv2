package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer(httpServer *http.Server) *Server {
	return &Server{HttpServer: httpServer}
}

func (s Server) Run() {
	err := s.HttpServer.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func (s Server) Stop(ctx context.Context) error {
	time.Sleep(time.Second * 2)
	err := s.HttpServer.Shutdown(ctx)
	log.Println("Stop http server")
	return err
}
