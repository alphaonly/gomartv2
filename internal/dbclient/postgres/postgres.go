package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alphaonly/gomartv2/internal/common"
	"github.com/alphaonly/gomartv2/internal/dbclient"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresClient struct {
	dataBaseURL string
	pool        *pgxpool.Pool
	conn        *pgxpool.Conn
}

func (pc postgresClient) GetConn() (*pgxpool.Conn, error) {
	if pc.conn == nil {
		return nil, fmt.Errorf(Message[8])
	}
	return pc.conn, nil
}

func (pc postgresClient) GetPull() (*pgxpool.Pool, error) {
	if pc.pool == nil {
		return nil, fmt.Errorf(Message[9])
	}
	return pc.pool, nil
}

func NewPostgresClient(ctx context.Context, dataBaseURL string) dbclient.DBClient {
	//get params
	s := postgresClient{dataBaseURL: dataBaseURL}
	//connect db
	var err error
	//s.conn, err = pgx.Connect(ctx, s.dataBaseURL)
	s.pool, err = pgxpool.New(ctx, s.dataBaseURL)
	if err != nil {
		common.LogFatalf(Message[0], err)
		return nil
	}
	// check users table exists
	err = CreateTable(ctx, s, checkIfUsersTableExists, createUsersTable)
	common.LogFatalf("error:", err)
	// check orders table exists
	err = CreateTable(ctx, s, checkIfOrdersTableExists, createOrdersTable)
	common.LogFatalf("error:", err)
	// check withdrawals table exists
	err = CreateTable(ctx, s, checkIfWithdrawalsTableExists, createWithdrawalsTable)
	common.LogFatalf("error:", err)

	return &s
}

func (pc *postgresClient) Connect(ctx context.Context) (ok bool) {
	ok = false
	var err error

	if pc.pool == nil {
		pc.pool, err = pgxpool.New(ctx, pc.dataBaseURL)
		common.LogFatalf(Message[0], err)
	}
	for i := 0; i < 10; i++ {
		pc.conn, err = pc.pool.Acquire(ctx)
		if err != nil {
			log.Println(Message[12] + " " + err.Error())
			time.Sleep(time.Millisecond * 200)
			continue
		}
		break
	}

	err = pc.conn.Ping(ctx)
	if err != nil {
		common.LogFatalf(Message[0], err)
	}

	ok = true
	return ok
}
