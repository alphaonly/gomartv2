package order

import (
	"context"
	"fmt"
	"strconv"

	"github.com/theplant/luhn"
)

var (
	ErrUserIsEmpty       = fmt.Errorf("400 user is empty")
	ErrBadOrderNumber    = fmt.Errorf("400 bad  number")
	ErrNoOrders          = fmt.Errorf("204 no orders")
	ErrNoLuhnNumber      = fmt.Errorf("422 not Lihn number")
	ErrOrderNumberExists = fmt.Errorf("200 order exists with user")
	ErrAnotherUsersOrder = fmt.Errorf("409 order exists with another user")
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
		ErrUserIsEmpty = fmt.Errorf("400 user is empty %v (%w)", userName, ErrUserIsEmpty)
		return nil, ErrUserIsEmpty
	}
	//getOrders
	orderslist, err := sr.Storage.GetOrdersList(ctx, userName)
	if err != nil {
		ErrNoOrders = fmt.Errorf(err.Error()+"(%w)", ErrNoOrders)
		ErrNoOrders = fmt.Errorf("204 no orders for user %v %w", userName, ErrNoOrders)
		return nil, ErrNoOrders
	}
	return orderslist, nil
}

func (sr service) ValidateOrderNumber(ctx context.Context, orderNumberStr string, user string) (orderNum int64, err error) {

	orderNumber, err := strconv.Atoi(orderNumberStr)
	if err != nil {

		ErrBadOrderNumber = fmt.Errorf(err.Error()+"(%w)", ErrBadOrderNumber)
		ErrBadOrderNumber = fmt.Errorf("400 order number bad number value %w", ErrBadOrderNumber)
		return 0, ErrBadOrderNumber
	}
	// order number format check
	if orderNumber <= 0 {
		ErrBadOrderNumber = fmt.Errorf("400 no order number zero or less(%w)", ErrBadOrderNumber)
		return int64(orderNumber), ErrBadOrderNumber
	}
	// orderNumber number validation according Luhn algorithm
	if !luhn.Valid(orderNumber) {
		ErrNoLuhnNumber = fmt.Errorf("422 no order number with Luhn: %v(%w)", orderNumber, ErrNoLuhnNumber)
		return int64(orderNumber), ErrNoLuhnNumber
	}
	// Check if orderNumber had already existed
	orderChk, err := sr.Storage.GetOrder(ctx, int64(orderNumber))
	if err != nil {
		return int64(orderNumber), nil
	}
	//Order exists, check user
	if user == orderChk.User {
		ErrOrderNumberExists = fmt.Errorf("200 order %v exists with user %v(%w)", orderNumber, user, ErrOrderNumberExists)
		return int64(orderNumber), ErrOrderNumberExists
	}
	ErrAnotherUsersOrder = fmt.Errorf("409 order %v exists with another user %v(%w)", orderNumber, orderChk.User, ErrAnotherUsersOrder)
	return int64(orderNumber), ErrAnotherUsersOrder
}
