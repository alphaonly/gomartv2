package mocks

import (
	"context"
	"fmt"
	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/schema"
)

type service struct {
}

func NewOrderStorage() order.Storage { return nil }

func NewService() (sr order.Service) {
	return service{}
}

var testUser = "testuser"

var testOrder1122 = order.Order{
	Order:   "1122",
	User:    testUser,
	Status:  order.NewOrder.Text,
	Accrual: 0,
	Created: schema.CreatedTime{},
}

func (sr service) GetUsersOrders(ctx context.Context, userName string) (orders order.Orders, err error) {
	// data validation
	if userName == "" {
		order.ErrUserIsEmpty = fmt.Errorf("400 user is empty %v (%w)", userName, order.ErrUserIsEmpty)
		return nil, order.ErrUserIsEmpty
	}
	//getOrders
	if userName == testUser {
		return order.Orders{1122: testOrder1122}, nil
	}
	return nil, order.ErrNoOrders

}

func (sr service) ValidateOrderNumber(ctx context.Context, orderNumberStr string, user string) (orderNum int64, err error) {

	return 0, err
}
