package user

import (
	"context"
	"errors"
	"fmt"
	"log"
)

// Constants to describe typical errors during manipulation with order entity
var (
	ErrUserPassEmpty  = fmt.Errorf("400 user or password is empty")
	ErrInternal       = fmt.Errorf("500 internal error: ")
	ErrLoginOccupied  = fmt.Errorf("409 login is occupied")
	ErrSaveUser       = fmt.Errorf("500 cannot save user in storage")
	ErrLogPassUnknown = errors.New("401 login or password is unknown")
)

// Service - an interface that implements the logic of manipulation with user entity
type Service interface {
	RegisterUser(ctx context.Context, u *User) (err error)                                         // Registration of a new user
	AuthenticateUser(ctx context.Context, u *User) (err error)                                     // Authentication of a user
	CheckIfUserAuthorized(ctx context.Context, login string, password string) (ok bool, err error) // checks if user currently authorized
	GetUserBalance(ctx context.Context, userName string) (response *BalanceResponseDTO, err error) // Gets balance of authorized user
}

type service struct {
	Storage Storage
}

// NewService - a factory that return the implementation of Service for user entity
func NewService(s Storage) (sr Service) {
	return &service{Storage: s}
}

// RegisterUser - implements logic of registration of a new user
func (sr service) RegisterUser(ctx context.Context, u *User) (err error) {
	// data validation
	if u.User == "" || u.Password == "" {
		return ErrUserPassEmpty
	}
	// Check if username exists
	userChk, err := sr.Storage.GetUser(ctx, u.User)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", u.User, ErrInternal)
		ErrInternal = fmt.Errorf("500 internal error in getting user %v: %w", u.User, ErrInternal)
		return ErrInternal
	}
	if userChk != nil {
		//login has already been occupied
		ErrLoginOccupied = fmt.Errorf("409 login %v is occupied (%w)", userChk.User, ErrLoginOccupied)
		return ErrLoginOccupied
	}
	err = sr.Storage.SaveUser(ctx, u)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", ErrInternal)
		ErrInternal = fmt.Errorf(" 500 cannot save user in storage %w", ErrInternal)
		return ErrInternal
	}
	return nil
}

// AuthenticateUser - implements logic of authentication of user
func (sr service) AuthenticateUser(ctx context.Context, u *User) (err error) {
	// data validation
	if u.User == "" || u.Password == "" {
		return ErrUserPassEmpty
	}
	// Check if username exists
	userInStorage, err := sr.Storage.GetUser(ctx, u.User)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", u.User, ErrInternal)
		ErrInternal = fmt.Errorf("500 internal error in getting user %v: %w", u.User, ErrInternal)
		log.Println(ErrInternal)
		return ErrInternal
	}
	if !u.Equals(userInStorage) {
		return ErrLogPassUnknown
	}

	return nil
}

// CheckIfUserAuthorized -  implements logic of check whether user is authorized
func (sr service) CheckIfUserAuthorized(ctx context.Context, login string, password string) (ok bool, err error) {
	// data validation
	if login == "" || password == "" {
		return false, ErrUserPassEmpty
	}
	// Check if username authorized
	u, err := sr.Storage.GetUser(ctx, login)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", ErrInternal)
		ErrInternal = fmt.Errorf("500 checking user authorization, can not get user from storage: %w", ErrInternal)
		log.Println(ErrInternal)
		return false, ErrInternal
	}
	if u == nil {
		ErrLogPassUnknown = fmt.Errorf("401 no user in storage means not authorized(%w)", ErrLogPassUnknown)
		return false, ErrLogPassUnknown
	}
	if !u.Equals(&User{User: login, Password: password}) {
		return false, nil
	}

	return true, nil
}

// GetUserBalance - implements logic of getting user's balance
func (sr service) GetUserBalance(ctx context.Context, userName string) (response *BalanceResponseDTO, err error) {
	// data validation
	if userName == "" {
		ErrUserPassEmpty = fmt.Errorf("400 username %v is empty(%w)", userName, ErrUserPassEmpty)
		return nil, ErrUserPassEmpty
	}
	//getUser
	user, err := sr.Storage.GetUser(ctx, userName)
	if err != nil {
		ErrInternal = fmt.Errorf(err.Error()+"(%w)", ErrInternal)
		ErrInternal = fmt.Errorf("500 unable to get user %v balance:  %w", userName, ErrInternal)
		return nil, ErrInternal
	}
	if user == nil {
		ErrLogPassUnknown = fmt.Errorf("401 unable to get user %v balance, as no user in storage(%w)", userName, ErrLogPassUnknown)
		return nil, ErrLogPassUnknown
	}
	return &BalanceResponseDTO{user.Accrual, user.Withdrawal}, nil
}
