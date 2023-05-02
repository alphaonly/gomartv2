package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alphaonly/gomartv2/internal/pkg/common/logging"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient/postgres"
	"log"
	"strings"

	"github.com/alphaonly/gomartv2/internal/domain/user"
)

type userStorage struct {
	client dbclient.DBClient
}

func NewStorage(client dbclient.DBClient) user.Storage {
	return &userStorage{client: client}

}

func (s userStorage) GetUser(ctx context.Context, name string) (u *user.User, err error) {
	if !s.client.Connect(ctx) {
		return nil, errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	d := DBUsersDTO{userID: sql.NullString{String: name, Valid: true}}
	row := conn.QueryRow(ctx, selectLineUsersTable, &d.userID)
	err = row.Scan(&d.userID, &d.password, &d.accrual, &d.withdrawal)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, nil
		}
		return nil, err
	}
	return &user.User{
		User:       d.userID.String,
		Password:   d.password.String,
		Accrual:    d.accrual.Float64,
		Withdrawal: d.withdrawal.Float64,
	}, nil
}

func (s userStorage) SaveUser(ctx context.Context, u *user.User) (err error) {
	if !s.client.Connect(ctx) {
		return errors.New(postgres.Message[0])
	}
	conn, err := s.client.GetConn()
	if err != nil {
		return err
	}
	defer conn.Release()

	d := DBUsersDTO{
		userID:     sql.NullString{String: u.User, Valid: true},
		password:   sql.NullString{String: u.Password, Valid: true},
		accrual:    sql.NullFloat64{Float64: u.Accrual, Valid: true},
		withdrawal: sql.NullFloat64{Float64: u.Withdrawal, Valid: true},
	}

	tag, err := conn.Exec(ctx, createOrUpdateIfExistsUsersTable, d.userID, d.password, d.accrual, d.withdrawal)
	logging.LogFatalf(postgres.Message[3], err)
	log.Println(tag)
	return err
}
