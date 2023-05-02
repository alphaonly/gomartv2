package user

import (
	"context"
)

type Storage interface {
	GetUser(ctx context.Context, name string) (u *User, err error)
	SaveUser(ctx context.Context, u *User) (err error)
}
