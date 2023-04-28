package user_test

import (
	"context"
	"log"
	"testing"

	"github.com/alphaonly/gomartv2/internal/domain/user"
	mockUser "github.com/alphaonly/gomartv2/internal/mocks/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mockUser.NewMockStorage(ctrl)

	tests := []struct {
		name    string
		regUser *user.User
		getUser *user.User
		setUser *user.User
		getErr  error
		setErr  error
		wantErr error
	}{
		{
			name:    "#1 Positive",
			regUser: &user.User{User: "testuser", Password: "password"},
			getUser: nil,
			getErr:  nil,
			setUser: &user.User{User: "testuser", Password: "password"},
			setErr:  nil,
			wantErr: nil,
		},
		{
			name:    "#2 Negative - no orders for user",
			regUser: &user.User{User: "testuser", Password: "password"},
			getUser: &user.User{User: "testuser", Password: "password"},
			getErr:  nil,
			setUser: nil,
			setErr:  user.ErrLoginOccupied,
			wantErr: user.ErrLoginOccupied,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			switch i {

			case 1:
				{
					s.EXPECT().GetUser(context.Background(), tt.regUser.User).Return(tt.getUser, tt.getErr)
				}
			default:
				{
					s.EXPECT().GetUser(context.Background(), tt.regUser.User).Return(tt.getUser, tt.getErr)
					s.EXPECT().SaveUser(context.Background(), tt.setUser).Return(tt.setErr)
				}
			}

			service := user.NewService(s)

			err := service.RegisterUser(context.Background(), tt.regUser)

			log.Println(err)

			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("Error %v but want %v", err, tt.wantErr)
			}

		})

	}
}
