package dbclient

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBClient interface {
	Connect(ctx context.Context) (ok bool)
	GetPull() (*pgxpool.Pool, error)
	GetConn() (*pgxpool.Conn, error)
}
