package composites

import (
	orderh "github.com/alphaonly/gomartv2/internal/adapters/api/order"
	orderd "github.com/alphaonly/gomartv2/internal/adapters/db/order"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/domain/user"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient"
)

// OrderComposite - a composite structure for orders
type OrderComposite struct {
	Storage order.Storage
	Service order.Service
	Handler orderh.Handler
}

// NewOrderComposite - it is a factory that returns an instance of order composite
func NewOrderComposite(dbClient dbclient.DBClient, userService user.Service, configuration *configuration.ServerConfiguration) *OrderComposite {
	storage := orderd.NewStorage(dbClient)
	service := order.NewService(storage)
	handler := orderh.NewHandler(storage, service, userService, configuration)
	return &OrderComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}
