package order

import "context"

// Storage - an interface that implements the logic for order data manipulation
type Storage interface {
	GetOrder(ctx context.Context, orderNumber int64) (o *Order, err error)     // Gets order data from a storage
	SaveOrder(ctx context.Context, o Order) (err error)                        // saves order data to a storage
	GetOrdersList(ctx context.Context, userName string) (ol Orders, err error) //gets list of all orders by given user
	GetNewOrdersList(ctx context.Context) (ol Orders, err error)               //gets list of new orders by given user
}
