package composites

import (
	userh "github.com/alphaonly/gomartv2/internal/adapters/api/user"
	userd "github.com/alphaonly/gomartv2/internal/adapters/db/user"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/user"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient"
)

// UserComposite - a composite structure for users
type UserComposite struct {
	Storage user.Storage
	Service user.Service
	Handler userh.Handler
}

// NewUserComposite - it is a factory that returns an instance of user composite
func NewUserComposite(dbClient dbclient.DBClient, configuration *configuration.ServerConfiguration) *UserComposite {
	storage := userd.NewStorage(dbClient)
	service := user.NewService(storage)
	handler := userh.NewHandler(storage, service, configuration)
	return &UserComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}
