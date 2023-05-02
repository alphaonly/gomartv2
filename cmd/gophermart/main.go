package main

import (
	"context"
	"github.com/alphaonly/gomartv2/internal/adapters/accrual"
	"github.com/alphaonly/gomartv2/internal/adapters/api"
	"github.com/alphaonly/gomartv2/internal/adapters/api/router"
	"github.com/alphaonly/gomartv2/internal/composites"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/pkg/common/logging"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient/postgres"
	"github.com/alphaonly/gomartv2/internal/pkg/server"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	dbclient := postgres.NewPostgresClient(ctx, cfg.DatabaseURI)

	UserComposite := composites.NewUserComposite(dbclient, cfg)
	OrderComposite := composites.NewOrderComposite(dbclient, cfg)
	WithdrawalComposite := composites.NewWithdrawalComposite(dbclient, cfg, UserComposite.Storage, OrderComposite.Service)

	handlerComposite := composites.NewHandlerComposite(
		api.NewHandler(cfg),
		UserComposite.Handler,
		OrderComposite.Handler,
		WithdrawalComposite.Handler,
	)

	// маршрутизация запросов обработчику
	rtr := router.NewRouter(handlerComposite)

	httpServer := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: rtr,
	}

	srv := server.NewServer(httpServer)
	acr := accrual.NewAccrual(cfg, OrderComposite.Storage)

	go srv.Run()
	go acr.Run(ctx)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	<-osSignal
	err := srv.Stop(ctx)
	logging.LogFatal(err)

}
