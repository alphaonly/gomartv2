package mocks

import (
	"context"
	"fmt"
	"github.com/alphaonly/gomartv2/internal/domain/withdrawal"
)

var (
	TestUser200 = "testuser200"
	TestUser500 = "testuser500"
	TestUser402 = "testuser402"
	TestUser422 = "testuser422"
	TestJSON    = []byte(fmt.Sprintf(`{"order":"%v","sum":%v}`, "2377225624", 751))
)

func NewWithdrawalStorage() withdrawal.Storage { return nil }

func NewService() (sr withdrawal.Service) {
	return &service{}
}

type service struct {
}

func (sr service) MakeUserWithdrawal(ctx context.Context, userName string, request withdrawal.UserWithdrawalRequestDTO) (err error) {
	// data validation
	if userName == "" {
		return withdrawal.ErrNoUser
	}

	if userName == TestUser200 {
		return nil
	}

	if userName == TestUser500 {
		return withdrawal.ErrInternal
	}

	if userName == TestUser402 {
		return withdrawal.ErrNoFunds
	}

	if userName == TestUser422 {
		return withdrawal.ErrOrderInvalid
	}
	return withdrawal.ErrInternal
}

func (sr service) GetUsersWithdrawals(ctx context.Context, userName string) (withdrawals *withdrawal.Withdrawals, err error) {
	return nil, nil
}
