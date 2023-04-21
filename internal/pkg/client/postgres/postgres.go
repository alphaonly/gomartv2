package postgres

import (
	"context"
	"log"
	"time"

	"github.com/alphaonly/gomartv2/internal/pkg/client"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresClient struct {
	dataBaseURL string
	pool        *pgxpool.Pool
	conn        *pgxpool.Conn
}

func NewPostgresClient(ctx context.Context, dataBaseURL string) client.DBClient {
	//get params
	s := postgresClient{dataBaseURL: dataBaseURL}
	//connect db
	var err error
	//s.conn, err = pgx.Connect(ctx, s.dataBaseURL)
	s.pool, err = pgxpool.New(ctx, s.dataBaseURL)
	if err != nil {
		logFatalf(message[0], err)
		return nil
	}
	// check users table exists
	err = createTable(ctx, s, checkIfUsersTableExists, createUsersTable)
	logFatalf("error:", err)
	// check orders table exists
	err = createTable(ctx, s, checkIfOrdersTableExists, createOrdersTable)
	logFatalf("error:", err)
	// check withdrawals table exists
	err = createTable(ctx, s, checkIfWithdrawalsTableExists, createWithdrawalsTable)
	logFatalf("error:", err)

	return &s
}

func (s *postgresClient) connectDB(ctx context.Context) (ok bool) {
	ok = false
	var err error

	if s.pool == nil {
		s.pool, err = pgxpool.New(ctx, s.dataBaseURL)
		client.LogFatalf(message[0], err)
	}
	for i := 0; i < 10; i++ {
		s.conn, err = s.pool.Acquire(ctx)
		if err != nil {
			log.Println(message[12] + " " + err.Error())
			time.Sleep(time.Millisecond * 200)
			continue
		}
		break
	}

	err = s.conn.Ping(ctx)
	if err != nil {
		logFatalf(message[0], err)
	}

	ok = true
	return ok
}
