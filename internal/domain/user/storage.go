package user

import (
	"context"
)

// Storage - an interface that implements the logic for user data manipulation
type Storage interface {
	GetUser(ctx context.Context, name string) (u *User, err error) //gets user data by user name
	SaveUser(ctx context.Context, u *User) (err error)             // saves user data
}
