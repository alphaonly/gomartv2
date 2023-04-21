package client

import (
	"context"
)

type DBClient interface {
	ConnectDB(ctx context.Context) (ok bool)
}
