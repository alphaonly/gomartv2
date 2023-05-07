package withdrawal

import (
	"context"
)

// Storage - an interface that implements the logic for withdrawal data manipulation
type Storage interface {
	SaveWithdrawal(ctx context.Context, w Withdrawal) (err error)                         // saves withdrawal data
	GetWithdrawalsList(ctx context.Context, userName string) (wl *Withdrawals, err error) // gets withdrawal data list by username
}
