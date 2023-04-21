package user

import (
	"context"
	"errors"
	"fmt"
	"log"
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
		return errors.New("400 user or password is empty")
	}
	// Check if username exists
	userChk, err := sr.Storage.GetUser(ctx, u.User)
	if err != nil {
		return fmt.Errorf("500 internal error in getting user %v: %w", u.User, err)
	}
	if userChk != nil {
		//login has already been occupied
		return errors.New("409 login " + userChk.User + " is occupied")
	}
	err = sr.Storage.SaveUser(ctx, u)
	if err != nil {
		return fmt.Errorf("cannot save user in storage %w", err)
	}
	return nil
}

func (sr service) AuthenticateUser(ctx context.Context, u *User) (err error) {
	// data validation
	if u.User == "" || u.Password == "" {
		return errors.New("400 user or password is empty")
	}
	// Check if username exists
	userInStorage, err := sr.Storage.GetUser(ctx, u.User)
	if err != nil {
		log.Printf("500 can not get user from storage:%v", err.Error())
		return fmt.Errorf("500 internal error in getting user %v: %w", u.User, err)
	}
	if !u.Equals(userInStorage) {
		return errors.New("401 login or password is unknown")
	}
	// sr.AuthorizedUsers[u.User] = true

	return nil
}

func (sr service) CheckIfUserAuthorized(ctx context.Context, login string, password string) (ok bool, err error) {
	// data validation
	if login == "" || password == "" {
		return false, errors.New("400 login or password is empty")
	}
	// Check if username authorized
	u, err := sr.Storage.GetUser(ctx, login)
	if err != nil {
		log.Printf("500 checking user authorization, can not get user from storage:%v", err.Error())
		return false, fmt.Errorf("500 get user internal error: %w", err)
	}
	if u == nil {
		return false, fmt.Errorf("401 no user in storage means not authorized: %w", err)
	}
	if !u.Equals(&User{User: login, Password: password}) {
		return false, nil
	}

	return true, nil
}

func (sr service) GetUserBalance(ctx context.Context, userName string) (response *UserBalanceResponse, err error) {
	// data validation
	if userName == "" {
		return nil, fmt.Errorf("400 username %v is empty", userName)
	}
	//getUser
	user, err := sr.Storage.GetUser(ctx, userName)
	if err != nil {
		return nil, fmt.Errorf("500 unable to get user %v balance:  %w", userName, err)
	}
	if user == nil {
		return nil, fmt.Errorf("401 unable to get user %v balance, as no user in storage", userName)
	}
	return &UserBalanceResponse{user.Accrual, user.Withdrawal}, nil
}
