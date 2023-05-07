// Package mocks - mocks various interfaces of user entity for testing

package mocks

import (
	"context"
	"github.com/alphaonly/gomartv2/internal/domain/user"
)

type service struct {
}

var (
	TestUser200 = "testuser200"
	TestPass200 = "password200"
	TestUser500 = "testuser500"
	TestPass500 = "password500"
	TestUser409 = "testuser409"
	TestPass409 = "password409"
)

func NewUserStorage() user.Storage { return nil }

func NewService() (sr user.Service) {
	return &service{}
}

func (sr service) RegisterUser(ctx context.Context, u *user.User) (err error) {
	// data validation
	if u.User == "" || u.Password == "" {
		return user.ErrUserPassEmpty
	}

	if u.User == TestUser200 && u.Password == TestPass200 {
		return nil
	}

	if u.User == TestUser500 && u.Password == TestPass500 {
		return user.ErrInternal
	}

	if u.User == TestUser409 && u.Password == TestPass409 {
		return user.ErrLoginOccupied
	}
	return user.ErrInternal
}

func (sr service) AuthenticateUser(ctx context.Context, u *user.User) (err error) {
	return nil
}

func (sr service) CheckIfUserAuthorized(ctx context.Context, login string, password string) (ok bool, err error) {
	return true, nil
}

func (sr service) GetUserBalance(ctx context.Context, userName string) (response *user.BalanceResponseDTO, err error) {
	return &user.BalanceResponseDTO{Current: 0, Withdrawn: 0}, nil
}
