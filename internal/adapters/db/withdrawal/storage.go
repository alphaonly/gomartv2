package withdrawal

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alphaonly/gomartv2/internal/pkg/common/logging"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient/postgres"
	"log"
	"time"

	"github.com/alphaonly/gomartv2/internal/domain/withdrawal"
	"github.com/alphaonly/gomartv2/internal/schema"
)

type withdrawalStorage struct {
	client dbclient.DBClient
}

// NewStorage - it is a factory that returns an instance of withdrawal's Storage implementation.
func NewStorage(client dbclient.DBClient) withdrawal.Storage {
	return &withdrawalStorage{client: client}

}

// SaveWithdrawal - an implementation of the function that saves withdrawal's data to postgres database
func (s withdrawalStorage) SaveWithdrawal(ctx context.Context, w withdrawal.Withdrawal) (err error) {

	if !s.client.Connect(ctx) {
		return errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	logging.LogFatalf("", err)
	defer conn.Release()

	w.Processed = schema.CreatedTime(time.Now())

	d := DBWithdrawalsDTO{
		userID:     sql.NullString{String: w.User, Valid: true},
		createdAt:  sql.NullString{String: time.Time(w.Processed).Format(time.RFC3339), Valid: true},
		orderID:    sql.NullString{String: w.Order, Valid: true},
		withdrawal: sql.NullFloat64{Float64: w.Withdrawal, Valid: true},
	}
	tag, err := conn.Exec(ctx, createOrUpdateIfExistsWithdrawalsTable, &d.userID, &d.createdAt, &d.orderID, &d.withdrawal)
	logging.LogFatalf(postgres.Message[7], err)
	log.Println(tag)
	return err
}

// GetWithdrawalsList - an implementation of the function that returns a user's list of withdrawals data from postgres database
func (s withdrawalStorage) GetWithdrawalsList(ctx context.Context, username string) (wl *withdrawal.Withdrawals, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	logging.LogFatalf("", err)

	defer conn.Release()

	wl = new(withdrawal.Withdrawals)

	d := DBWithdrawalsDTO{userID: sql.NullString{String: username, Valid: true}}

	rows, err := conn.Query(ctx, selectAllWithdrawalsTableByUser, &d.userID)
	if err != nil {
		log.Printf(postgres.Message[4], err)
		return nil, err
	}
	log.Printf("getting withdrawals for user %v", d.userID)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&d.userID, &d.createdAt, &d.orderID, &d.withdrawal)
		logging.LogFatalf(postgres.Message[5], err)
		created, err := time.Parse(time.RFC3339, d.createdAt.String)
		logging.LogFatalf(postgres.Message[6], err)
		log.Printf("got withdrawal for user %v: %v", d.userID, d)

		w := withdrawal.Withdrawal{
			User:       d.userID.String,
			Processed:  schema.CreatedTime(created),
			Order:      d.orderID.String,
			Withdrawal: d.withdrawal.Float64,
		}
		log.Printf("append  withdrawal to return list  : %v", w)
		*wl = append(*wl, w)
	}

	return wl, nil
}
