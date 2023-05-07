// Package server - http server functionality
package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server - http server cover
type Server struct {
	HTTPServer *http.Server
}

// NewServer - a factory that return pointer to a Server instance
func NewServer(httpServer *http.Server) *Server {
	return &Server{HTTPServer: httpServer}
}

// Run - runs http server listening of application
func (s Server) Run() {
	err := s.HTTPServer.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

// Stop - stops listening, shuts down http server
func (s Server) Stop(ctx context.Context) error {
	time.Sleep(time.Second * 2)
	err := s.HTTPServer.Shutdown(ctx)
	log.Println("Stop http server")
	return err
}
