package order_test

import (
	"context"
	"errors"

	"log"
	"testing"

	"github.com/alphaonly/gomartv2/internal/domain/order"
	mockOrder "github.com/alphaonly/gomartv2/internal/mocks/order"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mockOrder.NewMockStorage(ctrl)

	tests := []struct {
		name         string
		userName     string
		returnOrders order.Orders
		returnErr    error
		want         error
	}{
		{
			name:         "#1 Positive",
			userName:     "testuser",
			returnOrders: order.Orders{1233: order.Order{Order: "1233", User: "testuser", Status: order.NewOrder.Text}},
			returnErr:    nil,
			want:         nil,
		},
		{
			name:     "#2 Negative - no orders for user",
			userName: "testuser",
			// returnOrders: order.Orders{1233: order.Order{Order: "1233", User: "testuser2", Status: order.NewOrder.Text}},
			returnErr: order.ErrNoOrders,
			want:      order.ErrNoOrders,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {

			s.EXPECT().GetOrdersList(context.Background(), tt.userName).Return(tt.returnOrders, tt.returnErr)
			service := order.NewService(s)

			_, err := service.GetUsersOrders(context.Background(), tt.userName)
			log.Println(err)

			if !assert.Equal(t, true, errors.Is(err, tt.want)) {
				t.Errorf("Error %v but want %v", err, tt.want)
			}

		})

	}
}
