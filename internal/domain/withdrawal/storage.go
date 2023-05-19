package withdrawal

import (
	"context"
)

type Storage interface {
	SaveWithdrawal(ctx context.Context, w Withdrawal) (err error)
	GetWithdrawalsList(ctx context.Context, userName string) (wl *Withdrawals, err error)
}
