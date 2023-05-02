package composites

import (
	orderh "github.com/alphaonly/gomartv2/internal/adapters/api/order"
	orderd "github.com/alphaonly/gomartv2/internal/adapters/db/order"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient"
)

type OrderComposite struct {
	Storage order.Storage
	Service order.Service
	Handler orderh.Handler
}

func NewOrderComposite(dbClient dbclient.DBClient, configuration *configuration.ServerConfiguration) *OrderComposite {
	storage := orderd.NewStorage(dbClient)
	service := order.NewService(storage)
	handler := orderh.NewHandler(storage, service, configuration)
	return &OrderComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}
