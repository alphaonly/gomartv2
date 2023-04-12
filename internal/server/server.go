package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alphaonly/gomartv2/internal/server/accrual"

	conf "github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/server/handlers"
	stor "github.com/alphaonly/gomartv2/internal/server/storage/interfaces"
)

type Configuration struct {
	serverPort string
}

type Server struct {
	configuration   *conf.ServerConfiguration
	InternalStorage stor.Storage
	ExternalStorage stor.Storage
	handlers        *handlers.Handlers
	httpServer      *http.Server
	AccrualChecker  *accrual.Checker
}

func NewConfiguration(serverPort string) *Configuration {
	return &Configuration{serverPort: ":" + serverPort}
}

func New(
	configuration *conf.ServerConfiguration,
	ExStorage stor.Storage,
	handlers *handlers.Handlers,
	accrualChecker *accrual.Checker) (server Server) {
	return Server{
		configuration:   configuration,
		InternalStorage: handlers.Storage,
		ExternalStorage: ExStorage,
		handlers:        handlers,
		AccrualChecker:  accrualChecker,
	}
}

func (s Server) ListenData(ctx context.Context) {
	// err := http.ListenAndServe(s.configuration.Port, s.handlers.NewRouter())
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) Run(ctx context.Context) error {

	// маршрутизация запросов обработчику
	s.httpServer = &http.Server{
		Addr:    s.configuration.RunAddress,
		Handler: s.handlers.NewRouter(),
	}

	go s.ListenData(ctx)
	go s.AccrualChecker.Run(ctx)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	<-osSignal
	err := s.Shutdown(ctx)

	return err
}
func (s Server) Shutdown(ctx context.Context) error {
	time.Sleep(time.Second * 2)
	err := s.httpServer.Shutdown(ctx)
	log.Println("Server shutdown")
	return err
}
