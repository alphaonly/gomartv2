package order

import (
	"context"
	"fmt"
	"strconv"

	"github.com/theplant/luhn"
)

// Constants to describe typical errors during manipulation with order entity
var (
	ErrBadUserOrOrder    = fmt.Errorf("400 user is empty or bad order number")
	ErrNoOrders          = fmt.Errorf("204 no orders")
	ErrNoLuhnNumber      = fmt.Errorf("422 not Lihn number")
	ErrOrderNumberExists = fmt.Errorf("200 order exists with user")
	ErrAnotherUsersOrder = fmt.Errorf("409 order exists with another user")
)

// Service - an interface that implements the logic of manipulation with order entity
type Service interface {
	GetUsersOrders(ctx context.Context, userName string) (orders Orders, err error)                          // returns the list of orders for authorized user
	ValidateOrderNumber(ctx context.Context, orderNumberStr string, user string) (orderNum int64, err error) //Checks the inbound orders' number is valid
}

type service struct {
	Storage Storage
}

// NewService - a factory that return the implementation of Service for order entity
func NewService(s Storage) (sr Service) {
	return service{Storage: s}
}

// GetUsersOrders - implements logic of returning the list of orders for authorized user
func (sr service) GetUsersOrders(ctx context.Context, userName string) (orders Orders, err error) {
	// data validation
	if userName == "" {
		ErrBadUserOrOrder = fmt.Errorf("400 user is empty %v (%w)", userName, ErrBadUserOrOrder)
		return nil, ErrBadUserOrOrder
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

// ValidateOrderNumber - implements logic of the check for inbound orders' numbers
func (sr service) ValidateOrderNumber(ctx context.Context, orderNumberStr string, user string) (orderNum int64, err error) {

	orderNumber, err := strconv.Atoi(orderNumberStr)
	if err != nil {
		ErrBadUserOrOrder = fmt.Errorf(err.Error()+"(%w)", ErrBadUserOrOrder)
		ErrBadUserOrOrder = fmt.Errorf("400 order number bad number value %w", ErrBadUserOrOrder)
		return 0, ErrBadUserOrOrder
	}
	// order number format check
	if orderNumber <= 0 {
		ErrBadUserOrOrder = fmt.Errorf("400 no order number zero or less(%w)", ErrBadUserOrOrder)
		return int64(orderNumber), ErrBadUserOrOrder
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
