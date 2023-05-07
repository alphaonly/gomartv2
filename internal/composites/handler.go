// Package composites - this is a package that contains composites structures for services, handlers, storages of various entities
package composites

import (
	"github.com/alphaonly/gomartv2/internal/adapters/api"
	"github.com/alphaonly/gomartv2/internal/adapters/api/order"
	"github.com/alphaonly/gomartv2/internal/adapters/api/user"
	"github.com/alphaonly/gomartv2/internal/adapters/api/withdrawal"
)

// HandlerComposite - a composite structure for handlers
type HandlerComposite struct {
	Common     api.Handler
	User       user.Handler
	Order      order.Handler
	Withdrawal withdrawal.Handler
}

// NewHandlerComposite - it is a factory that returns an instance of handler composite
func NewHandlerComposite(
	common api.Handler,
	user user.Handler,
	order order.Handler,
	withdrawal withdrawal.Handler) *HandlerComposite {

	return &HandlerComposite{
		Common:     common,
		User:       user,
		Order:      order,
		Withdrawal: withdrawal,
	}
}
