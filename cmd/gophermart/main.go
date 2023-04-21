package main

func main() {

	// configuration := conf.NewServerConf(conf.UpdateSCFromEnvironment, conf.UpdateSCFromFlags)

	// var (
	// 	externalStorage stor.Keeper
	// 	internalStorage stor.Keeper
	// )

	// externalStorage = nil
	// internalStorage = db.NewDBStorage(context.Background(), configuration.DatabaseURI)

	// handlers := &handlers.Handlers{
	// 	Storage:       internalStorage,
	// 	Conf:          conf.ServerConfiguration{DatabaseURI: configuration.DatabaseURI},
	// 	EntityHandler: handlers.NewEntityHandler(internalStorage),
	// }
	// accrualChecker := accrual.NewChecker(configuration.AccrualSystemAddress, configuration.AccrualTime, internalStorage)

	// gmServer := server.New(configuration, externalStorage, handlers, accrualChecker)

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// err := gmServer.Run(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
