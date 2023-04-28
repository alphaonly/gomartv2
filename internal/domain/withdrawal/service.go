package withdrawal

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/domain/user"
)

var (
	ErrEmptyUser      = fmt.Errorf("400 user is empty")
	ErrOrderInvalid   = fmt.Errorf("422 order number invalid")
	ErrNoUser         = fmt.Errorf("401 no user")
	ErrNoFunds        = fmt.Errorf("402 no funds")
	ErrUserUpdate     = fmt.Errorf("500 user update error")
	ErrSaveWithdrawal = fmt.Errorf("500 safe withdrawal error")
	ErrGetWithdrawal  = fmt.Errorf("500 get withdrawal error")
	ErrNoWithdrawal   = fmt.Errorf("204 no withdrawals for user")
)

type Service interface {
	MakeUserWithdrawal(ctx context.Context, userName string, request UserWithdrawalRequestDTO) (err error)
	GetUsersWithdrawals(ctx context.Context, userName string) (withdrawals *Withdrawals, err error)
}

type service struct {
	Storage      Storage
	UStorage     user.Storage
	OrderService order.Service
}

func NewService(storage Storage, userStorage user.Storage, orderService order.Service) (sr Service) {
	return &service{
		Storage:      storage,
		UStorage:     userStorage,
		OrderService: orderService,
	}
}

func (sr service) MakeUserWithdrawal(ctx context.Context, userName string, ByUserRequestDTO UserWithdrawalRequestDTO) (err error) {
	// data validation
	if userName == "" {
		ErrEmptyUser = fmt.Errorf("400 user %v is empty", userName)
		return ErrEmptyUser
	}
	//check order number
	orderNumber, err := sr.OrderService.ValidateOrderNumber(ctx, ByUserRequestDTO.Order, userName)
	if err != nil {
		ErrOrderInvalid = fmt.Errorf("422 order number %v invalid %w", orderNumber, err)
		if !strings.Contains(errors.Unwrap(ErrOrderInvalid).Error(), "200") {
			return ErrOrderInvalid
		}
	}
	//getUser
	user, err := sr.UStorage.GetUser(ctx, userName)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("401 unable to make user withdrawal as no user %v in storage (%w)", userName, ErrNoUser)
	}
	// Calculate new accrual
	newAccrual := user.Accrual - ByUserRequestDTO.Sum
	if newAccrual < 0 {
		ErrNoFunds = fmt.Errorf("402 insufficient funds for withdrawal as user has %v and tries to withdraw %v (%w)", user.Accrual, ByUserRequestDTO.Sum, ErrNoFunds)
		return ErrNoFunds
	}
	//Update Users data
	user.Accrual = newAccrual
	user.Withdrawal += ByUserRequestDTO.Sum
	err = sr.UStorage.SaveUser(ctx, user)
	if err != nil {
		ErrUserUpdate = fmt.Errorf(err.Error()+"(%w)", ErrUserUpdate)
		ErrUserUpdate = fmt.Errorf("500 can not update data of user %v after withdrawal attempt on order %v %w", userName, orderNumber, ErrUserUpdate) 
		return 
	}
	//Add withdrawal
	w := Withdrawal{
		User: userName,
		// Processed:
		Order:      ByUserRequestDTO.Order,
		Withdrawal: ByUserRequestDTO.Sum,
	}
	err = sr.Storage.SaveWithdrawal(ctx, w)
	if err != nil {
		ErrSaveWithdrawal = fmt.Errorf(err.Error()+"(%w)", ErrGetWithdrawal)
		ErrSaveWithdrawal = fmt.Errorf("500 can not create withdrawal data for user %v after withdrawal attempt on order %v %w", userName, orderNumber, ErrGetWithdrawal)
		return ErrSaveWithdrawal
	}
	return nil
}

func (sr service) GetUsersWithdrawals(ctx context.Context, userName string) (withdrawals *Withdrawals, err error) {
	// data validation
	if userName == "" {
		ErrEmptyUser = fmt.Errorf("400 user %v is empty", userName)
		return nil, ErrEmptyUser
	}
	//getOrders
	wList, err := sr.Storage.GetWithdrawalsList(ctx, userName)
	if err != nil {
		ErrGetWithdrawal = fmt.Errorf(err.Error()+"(%w)", ErrGetWithdrawal)
		ErrGetWithdrawal = fmt.Errorf("500 internal error on getting withdrawals for user %v %w ", userName, ErrGetWithdrawal)
		return nil, ErrGetWithdrawal
	}
	if len(*wList) == 0 {
		ErrNoWithdrawal = fmt.Errorf("204 no withdrawals for user %v (%w)", userName, ErrNoWithdrawal)
		return nil, ErrNoWithdrawal

	}

	sort.Sort(ByTimeDescending(*wList))
	return wList, nil
}
