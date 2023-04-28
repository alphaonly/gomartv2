package postgres

import (
	"context"
	"fmt"
	"log"
	"reflect"
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
	// func (pc postgresClient) GetConn() (*pgxpool.Conn, error) {
	if reflect.ValueOf(pc.conn).IsNil() {
		// if pc.conn == nil {
		return nil, fmt.Errorf(Message[8])
	}
	return pc.conn, nil
}

func (pc postgresClient) GetPull() (*pgxpool.Pool, error) {
	if reflect.ValueOf(pc.pool).IsNil() {
		// if pc.pool == nil {
		return nil, fmt.Errorf(Message[9])
	}
	return pc.pool, nil
}

func NewPostgresClient(ctx context.Context, dataBaseURL string) dbclient.DBClient {
	//get params
	pc := postgresClient{dataBaseURL: dataBaseURL}
	//connect db
	var err error
	var p *pgxpool.Pool

	pc.pool, err = p, err
	// pgxpool.New(ctx, pc.dataBaseURL)
	if err != nil {
		common.LogFatalf(Message[0], err)
		return nil
	}

	return &pc
}

func (pc *postgresClient) Connect(ctx context.Context) (ok bool) {
	ok = false
	var err error

	if reflect.TypeOf(pc.pool) == nil {

		pc.pool, err = pgxpool.New(ctx, pc.dataBaseURL)
		common.LogFatalf(Message[0], err)
	}
	for i := 0; i < 10; i++ {
		pc.conn, err =  pc.pool.Acquire(ctx)

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

func (pc postgresClient) CheckTables(ctx context.Context) error {
	if reflect.TypeOf(pc.conn) == nil {
		return fmt.Errorf(Message[9])
	}
	var err error
	// check users table exists
	err = CreateTable(ctx, pc, checkIfUsersTableExists, createUsersTable)
	common.LogFatalf("error:", err)
	// check orders table exists
	err = CreateTable(ctx, pc, checkIfOrdersTableExists, createOrdersTable)
	common.LogFatalf("error:", err)
	// check withdrawals table exists
	err = CreateTable(ctx, pc, checkIfWithdrawalsTableExists, createWithdrawalsTable)
	common.LogFatalf("error:", err)

	return nil
}
