// Package dbclient - a database client for connection
package dbclient

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBClient - an interface that implements a database client for connection to it
type DBClient interface {
	Connect(ctx context.Context) (ok bool)
	GetPull() (*pgxpool.Pool, error)
	GetConn() (*pgxpool.Conn, error)
}
