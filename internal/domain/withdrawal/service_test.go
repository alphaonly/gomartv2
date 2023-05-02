package withdrawal_test

import (
	"context"
	"errors"
	"strconv"

	"testing"

	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/domain/user"
	"github.com/alphaonly/gomartv2/internal/domain/withdrawal"
	mockOrder "github.com/alphaonly/gomartv2/internal/mocks/order"
	mockUser "github.com/alphaonly/gomartv2/internal/mocks/user"
	mockWithdrawal "github.com/alphaonly/gomartv2/internal/mocks/withdrawal"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMakeUserWithdrawal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wdStor := mockWithdrawal.NewMockStorage(ctrl)
	orStor := mockOrder.NewMockStorage(ctrl)
	usStor := mockUser.NewMockStorage(ctrl)

	var (
		userName               = "testuser"
		orderNumbInt     int64 = 12559
		orderNumbStr           = strconv.FormatInt(orderNumbInt, 10)
		byUserRequestDTO       = withdrawal.UserWithdrawalRequestDTO{
			Order: orderNumbStr,
			Sum:   80,
		}
	)
	tests := []struct {
		name            string
		userNameRequest string

		ByUserRequestDTO withdrawal.UserWithdrawalRequestDTO

		orderNumber      int64
		getOrderResponse *order.Order
		getUserResponse  *user.User

		setUserRequest     *user.User
		setUserErrResponse error

		setWdRequest     *withdrawal.Withdrawal
		setWdErrResponse error
		wantErr          error
	}{
		{
			name: "#1 Positive",
			//MakeUserWithdrawal request data
			userNameRequest:  userName,
			ByUserRequestDTO: byUserRequestDTO,
			//GetOrder	response data
			getOrderResponse: &order.Order{
				Order:   orderNumbStr,
				User:    userName,
				Status:  order.ProcessedOrder.Text,
				Accrual: 100,
			},
			orderNumber: orderNumbInt,
			//GetUser response data
			getUserResponse: &user.User{
				User:       userName,
				Accrual:    100,
				Withdrawal: 0,
			},
			//SaveUser requestData
			setUserRequest: &user.User{
				User:       userName,
				Accrual:    20,
				Withdrawal: 80,
			},
			setUserErrResponse: nil,

			setWdRequest: &withdrawal.Withdrawal{
				User: userName,
				//Processed:  schema.CreatedTime(time.Now()),
				Order:      byUserRequestDTO.Order,
				Withdrawal: byUserRequestDTO.Sum,
			},
			setWdErrResponse: nil,
			wantErr:          nil,
		},
		{
			name: "#2 Negative - Accrual is insufficient to withdraw",
			//MakeUserWithdrawal request data
			userNameRequest:  userName,
			ByUserRequestDTO: byUserRequestDTO,
			//GetOrder	response data
			getOrderResponse: &order.Order{
				Order:   orderNumbStr,
				User:    userName,
				Status:  order.ProcessedOrder.Text,
				Accrual: 100,
			},
			orderNumber: orderNumbInt,
			//GetUser response data
			getUserResponse: &user.User{
				User:       userName,
				Accrual:    10,
				Withdrawal: 0,
			},
			//SaveUser requestData
			setUserRequest: &user.User{
				User:       userName,
				Accrual:    20,
				Withdrawal: 80,
			},
			setUserErrResponse: nil,

			setWdRequest: &withdrawal.Withdrawal{
				User: userName,
				//Processed:  schema.CreatedTime(time.Now()),
				Order:      byUserRequestDTO.Order,
				Withdrawal: byUserRequestDTO.Sum,
			},
			setWdErrResponse: withdrawal.ErrNoFunds,
			wantErr:          withdrawal.ErrNoFunds,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			switch i {

			case 0:
				{
					orStor.EXPECT().GetOrder(context.Background(), tt.orderNumber).Return(tt.getOrderResponse, nil)
					usStor.EXPECT().GetUser(context.Background(), tt.setWdRequest.User).Return(tt.getUserResponse, nil)
					usStor.EXPECT().SaveUser(context.Background(), tt.setUserRequest).Return(tt.setUserErrResponse)
					wdStor.EXPECT().SaveWithdrawal(context.Background(), *tt.setWdRequest).Return(tt.setWdErrResponse)
				}
			case 1:
				{
					orStor.EXPECT().GetOrder(context.Background(), tt.orderNumber).Return(tt.getOrderResponse, nil)
					usStor.EXPECT().GetUser(context.Background(), tt.setWdRequest.User).Return(tt.getUserResponse, nil)
				}

			}

			orderService := order.NewService(orStor)

			service := withdrawal.NewService(wdStor, usStor, orderService)

			err := service.MakeUserWithdrawal(context.Background(), tt.userNameRequest, tt.ByUserRequestDTO)

			t.Log(err)

			if !assert.Equal(t, true, errors.Is( err, tt.wantErr)) {
				t.Errorf("Error %v but want %v", err, tt.wantErr)
			}

		})

	}
}
