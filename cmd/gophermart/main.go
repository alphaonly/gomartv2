package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/alphaonly/gomartv2/internal/adapters/accrual"
	"github.com/alphaonly/gomartv2/internal/adapters/api"
	"github.com/alphaonly/gomartv2/internal/adapters/api/router"
	"github.com/alphaonly/gomartv2/internal/composites"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/pkg/common/logging"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient/postgres"
	"github.com/alphaonly/gomartv2/internal/pkg/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	dbclient := postgres.NewPostgresClient(ctx, cfg.DatabaseURI)

	userComposite := composites.NewUserComposite(dbclient, cfg)
	orderComposite := composites.NewOrderComposite(dbclient, userComposite.Service, cfg)
	withdrawalComposite := composites.NewWithdrawalComposite(dbclient, cfg, userComposite.Storage, orderComposite.Service)

	handlerComposite := composites.NewHandlerComposite(
		api.NewHandler(cfg),
		userComposite.Handler,
		orderComposite.Handler,
		withdrawalComposite.Handler,
	)

	// маршрутизация запросов обработчику
	rtr := router.NewRouter(handlerComposite)

	httpServer := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: rtr,
	}

	srv := server.NewServer(httpServer)
	acr := accrual.NewAccrual(cfg, orderComposite.Storage, userComposite.Storage)

	go srv.Run()
	go acr.Run(ctx)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	<-osSignal
	err := srv.Stop(ctx)
	logging.LogFatal(err)

}
