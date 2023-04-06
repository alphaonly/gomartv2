package main

import (
	"context"
	conf "github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/server"
	"github.com/alphaonly/gomartv2/internal/server/accrual"
	"github.com/alphaonly/gomartv2/internal/server/handlers"
	db "github.com/alphaonly/gomartv2/internal/server/storage/implementations/dbstorage"
	stor "github.com/alphaonly/gomartv2/internal/server/storage/interfaces"
	"log"
)

func main() {

	configuration := conf.NewServerConf(conf.UpdateSCFromEnvironment, conf.UpdateSCFromFlags)

	var (
		externalStorage stor.Storage
		internalStorage stor.Storage
	)

	externalStorage = nil
	internalStorage = db.NewDBStorage(context.Background(), configuration.DatabaseURI)

	handlers := &handlers.Handlers{
		Storage:       internalStorage,
		Conf:          conf.ServerConfiguration{DatabaseURI: configuration.DatabaseURI},
		EntityHandler: handlers.NewEntityHandler(internalStorage),
	}
	accrualChecker := accrual.NewChecker(configuration.AccrualSystemAddress, configuration.AccrualTime, internalStorage)

	gmServer := server.New(configuration, externalStorage, handlers, accrualChecker)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := gmServer.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
