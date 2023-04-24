package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alphaonly/gomartv2/internal/adapters/accrual"
	"github.com/alphaonly/gomartv2/internal/adapters/api"
	"github.com/alphaonly/gomartv2/internal/adapters/api/router"
	"github.com/alphaonly/gomartv2/internal/composites"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/dbclient/postgres"
)

func main() {

	configuration := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	dbclient := postgres.NewPostgresClient(context.Background(), configuration.DatabaseURI)

	UserComposite := composites.NewUserComposite(dbclient, configuration)
	OrderComposite := composites.NewOrderComposite(dbclient, configuration)
	WithdrawalComposite := composites.NewWithdrawalComposite(dbclient, configuration, UserComposite.Storage, OrderComposite.Service)

	handlerComposite := composites.NewHandlerComposite(
		api.NewHandler(configuration),
		UserComposite.Handler,
		OrderComposite.Handler,
		WithdrawalComposite.Handler,
	)

	// маршрутизация запросов обработчику
	httpServer := &http.Server{
		Addr:    configuration.RunAddress,
		Handler: router.NewRouter(handlerComposite),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ListenData(ctx,httpServer)
	go accrual.NewAccrual(configuration, 
		OrderComposite.Storage).Run(ctx)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	<-osSignal
	err := Shutdown(ctx,httpServer)
	if err!=nil{
		log.Fatal(err)
	}
}

func ListenData(ctx context.Context, httpServer *http.Server) {
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func Shutdown(ctx context.Context, httpServer *http.Server) error {
	time.Sleep(time.Second * 2)
	err := httpServer.Shutdown(ctx)
	log.Println("Shutdown http server")
	return err
}
