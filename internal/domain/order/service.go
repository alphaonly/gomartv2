package order

import (
	"context"
	"fmt"
	"strconv"

	"github.com/theplant/luhn"
)

type Service interface {
	GetUsersOrders(ctx context.Context, userName string) (orders Orders, err error)
	ValidateOrderNumber(ctx context.Context, orderNumberStr string, user string) (orderNum int64, err error)
}

type service struct {
	Storage Storage
}

func NewService(s Storage) (sr Service) {
	return service{Storage: s}
}

func (sr service) GetUsersOrders(ctx context.Context, userName string) (orders Orders, err error) {
	// data validation
	if userName == "" {
		return nil, fmt.Errorf("400 user %v is empty", userName)
	}
	//getOrders
	orderslist, err := sr.Storage.GetOrdersList(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("204 no orders for user %v %w", userName, err)
	}
	return orderslist, nil
}

func (sr service) ValidateOrderNumber(ctx context.Context, orderNumberStr string, user string) (orderNum int64, err error) {

	orderNumber, err := strconv.Atoi(orderNumberStr)
	if err != nil {
		return 0, fmt.Errorf("400 order number bad number value %w", err)
	}
	// order number format check
	if orderNumber <= 0 {
		return int64(orderNumber), fmt.Errorf("400 no order number zero or less:%v", orderNumber)
	}
	// orderNumber number validation according Luhn algorithm
	if !luhn.Valid(orderNumber) {
		return int64(orderNumber), fmt.Errorf("422 no order number with Luhn: %v", orderNumber)
	}
	// Check if orderNumber had already existed
	orderChk, err := sr.Storage.GetOrder(ctx, int64(orderNumber))
	if err != nil {
		return int64(orderNumber), nil
	}
	//Order exists, check user
	if user == orderChk.User {
		return int64(orderNumber), fmt.Errorf("200 order %v exists with user %v", orderNumber, user)
	}
	return int64(orderNumber), fmt.Errorf("409 order %v exists with another user %v", orderNumber, orderChk.User)
}
