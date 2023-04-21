package withdrawal

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/domain/user"
	"github.com/alphaonly/gomartv2/internal/schema"
)

type Service interface {
	MakeUserWithdrawal(ctx context.Context, userName string, request UserWithdrawalRequest) (err error)
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

func (sr service) MakeUserWithdrawal(ctx context.Context, userName string, request UserWithdrawalRequest) (err error) {
	// data validation
	if userName == "" {
		return fmt.Errorf("400 user %v is empty", userName)
	}
	//check order number
	orderNumber, err := sr.OrderService.ValidateOrderNumber(ctx, request.Order, userName)
	if err != nil {
		return fmt.Errorf("422 order number invalid %v %w", orderNumber, err)
	}
	//getUser
	user, err := sr.UStorage.GetUser(ctx, userName)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("401 unable to make user withdrawal as no user %v in storage", userName)
	}
	//Calculate new accrual
	newAccrual := user.Accrual - request.Sum
	if newAccrual < 0 {
		return fmt.Errorf("402 insufficient funds for withdrawal")
	}
	//Update Users data
	user.Accrual = newAccrual
	user.Withdrawal += request.Sum
	err = sr.UStorage.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("500 can not update data of user %v after withrawal attempt on order %v %w", userName, orderNumber, err)
	}
	//Add withdrawal
	w := Withdrawal{
		User:       userName,
		Processed:  schema.CreatedTime(time.Now()),
		Order:      request.Order,
		Withdrawal: request.Sum,
	}
	err = sr.Storage.SaveWithdrawal(ctx, w)
	if err != nil {
		return fmt.Errorf("500 can not create withdrawal data for user %v after withdrawal attempt on order %v %w", userName, orderNumber, err)
	}
	return nil
}
func (sr service) GetUsersWithdrawals(ctx context.Context, userName string) (withdrawals *Withdrawals, err error) {
	// data validation
	if userName == "" {
		return nil, fmt.Errorf("400 user %v is empty", userName)
	}
	//getOrders
	wList, err := sr.Storage.GetWithdrawalsList(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("500 internal error on getting withdrawals for user %v %w", userName, err)
	}
	if len(*wList) == 0 {
		return nil, fmt.Errorf("204 no withdrawals for user %v %w", userName, err)
	}

	sort.Sort(ByTimeDescending(*wList))
	return wList, nil
}
