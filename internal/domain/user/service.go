package user

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var (
	ErrUserPassEmpty  = fmt.Errorf("400 user or password is empty")
	ErrInternal       = fmt.Errorf("500 internal error: ")
	ErrLoginOccupied  = fmt.Errorf("409 login is occupied")
	ErrSaveUser       = fmt.Errorf("500 cannot save user in storage")
	ErrLogPassUnknown = errors.New("401 login or password is unknown")
)

type Service interface {
	RegisterUser(ctx context.Context, u *User) (err error)
	AuthenticateUser(ctx context.Context, u *User) (err error)
	CheckIfUserAuthorized(ctx context.Context, login string, password string) (ok bool, err error)
	GetUserBalance(ctx context.Context, userName string) (response *UserBalanceResponse, err error)
}

type service struct {
	Storage Storage
}

func NewService(s Storage) (sr Service) {
	return &service{Storage: s}
}

func (sr service) RegisterUser(ctx context.Context, u *User) (err error) {
	// data validation
	if u.User == "" || u.Password == "" {
		return ErrUserPassEmpty
	}
	// Check if username exists
	userChk, err := sr.Storage.GetUser(ctx, u.User)
	if err != nil {
		ErrInternal = fmt.Errorf("500 internal error in getting user %v: %w", u.User, err)
		return ErrInternal
	}
	if userChk != nil {
		//login has already been occupied
		ErrLoginOccupied = errors.New("409 login " + userChk.User + " is occupied")
		return ErrLoginOccupied
	}
	err = sr.Storage.SaveUser(ctx, u)
	if err != nil {
		ErrInternal = fmt.Errorf(" 500 cannot save user in storage %w", err)
		return ErrInternal
	}
	return nil
}

func (sr service) AuthenticateUser(ctx context.Context, u *User) (err error) {
	// data validation
	if u.User == "" || u.Password == "" {
		return ErrUserPassEmpty
	}
	// Check if username exists
	userInStorage, err := sr.Storage.GetUser(ctx, u.User)
	if err != nil {
		ErrInternal = fmt.Errorf("500 internal error in getting user %v: %w", u.User, err)
		log.Println(ErrInternal)
		return ErrInternal
	}
	if !u.Equals(userInStorage) {
		return ErrLogPassUnknown
	}

	return nil
}

func (sr service) CheckIfUserAuthorized(ctx context.Context, login string, password string) (ok bool, err error) {
	// data validation
	if login == "" || password == "" {
		return false, ErrUserPassEmpty
	}
	// Check if username authorized
	u, err := sr.Storage.GetUser(ctx, login)
	if err != nil {		
		ErrInternal = fmt.Errorf("500 checking user authorization, can not get user from storage: %w", err)
		log.Println(ErrInternal)
		return false, ErrInternal
	}
	if u == nil {
		ErrLogPassUnknown = fmt.Errorf("401 no user in storage means not authorized: %w", err)
		return false, ErrLogPassUnknown
	}
	if !u.Equals(&User{User: login, Password: password}) {
		return false, nil
	}

	return true, nil
}

func (sr service) GetUserBalance(ctx context.Context, userName string) (response *UserBalanceResponse, err error) {
	// data validation
	if userName == "" {
		ErrUserPassEmpty = fmt.Errorf("400 username %v is empty", userName)
		return nil, ErrUserPassEmpty
	}
	//getUser
	user, err := sr.Storage.GetUser(ctx, userName)
	if err != nil {
		ErrInternal = fmt.Errorf("500 unable to get user %v balance:  %w", userName, err)
		return nil, ErrInternal
	}
	if user == nil {
		ErrLogPassUnknown = fmt.Errorf("401 unable to get user %v balance, as no user in storage", userName)
		return nil, ErrLogPassUnknown
	}
	return &UserBalanceResponse{user.Accrual, user.Withdrawal}, nil
}
