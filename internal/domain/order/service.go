package order

import (
	"context"
	"fmt"
	"strconv"

	"github.com/theplant/luhn"
)
var (
	Err                  error
	ErrUserIsEmpty       = fmt.Errorf("400 user is empty %w", Err)
	ErrBadOrderNumber    = fmt.Errorf("400 order number bad number value %w", Err)
	ErrNoOrders          = fmt.Errorf("204 no orders for user %w", Err)
	ErrNoLuhnNumber      = fmt.Errorf("422 no order number with Luhn: %w", Err)
	ErrOrderNumberExists = fmt.Errorf("200 order exists with user %w", Err)
	ErrAnotherUsersOrder = fmt.Errorf("409 order exists with another user %w", Err)
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
		Err = fmt.Errorf("400 user is empty %v", userName)
		return nil, ErrUserIsEmpty
	}
	//getOrders
	orderslist, err := sr.Storage.GetOrdersList(ctx, userName)
	if err != nil {
		Err = fmt.Errorf("204 no orders for user %v %w", userName, err)
		return nil, ErrNoOrders
	}
	return orderslist, nil
}

func (sr service) ValidateOrderNumber(ctx context.Context, orderNumberStr string, user string) (orderNum int64, err error) {

	orderNumber, err := strconv.Atoi(orderNumberStr)
	if err != nil {
		Err = fmt.Errorf("400 order number bad number value %w", err)
		return 0, ErrBadOrderNumber
	}
	// order number format check
	if orderNumber <= 0 {
		Err = fmt.Errorf("400 no order number zero or less:%v", orderNumber)
		return int64(orderNumber), ErrBadOrderNumber
	}
	// orderNumber number validation according Luhn algorithm
	if !luhn.Valid(orderNumber) {
		Err = fmt.Errorf("422 no order number with Luhn: %v", orderNumber)
		return int64(orderNumber), ErrNoLuhnNumber
	}
	// Check if orderNumber had already existed
	orderChk, err := sr.Storage.GetOrder(ctx, int64(orderNumber))
	if err != nil {
		return int64(orderNumber), nil
	}
	//Order exists, check user
	if user == orderChk.User {
		Err = fmt.Errorf("200 order %v exists with user %v", orderNumber, user)
		return int64(orderNumber), ErrOrderNumberExists
	}
	Err = fmt.Errorf("409 order %v exists with another user %v", orderNumber, orderChk.User)
	return int64(orderNumber), ErrAnotherUsersOrder
}
