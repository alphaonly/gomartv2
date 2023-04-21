package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/alphaonly/gomartv2/internal/schema"
	"github.com/jackc/pgx/v5/pgxpool"
)


// SaveOrder(ctx context.Context, o schema.Order) (err error)
// GetOrdersList(ctx context.Context, u schema.User) (wl schema.Orders, err error)


// -d=postgres://postgres:mypassword@localhost:5432/yandex
const (
	selectLineOrdersTable           = `SELECT order_id, user_id, status, accrual, uploaded_at FROM public.orders WHERE order_id=$1;`
	selectAllOrdersTableByUser      = `SELECT order_id, user_id, status, accrual, uploaded_at FROM public.orders WHERE user_id = $1;`
	selectAllOrdersTableByStatus    = `SELECT order_id, user_id, status, accrual, uploaded_at  FROM public.orders WHERE status = $1;`

	createOrUpdateIfExistsOrdersTable = `
	  INSERT INTO public.orders (order_id, user_id, status,accrual,uploaded_at) 
	  VALUES ($1, $2, $3,$4, $5)
	  ON CONFLICT (order_id,user_id) DO UPDATE 
		SET status 		= $3,
		    accrual 	= $4,
			uploaded_at = $5; 
		`
	

	createOrdersTable = `create table public.orders
	(	order_id bigint not null, 
		user_id varchar(40) not null,
		status integer,		
		accrual double precision,
		uploaded_at TEXT not null, 
		primary key (order_id,user_id)
	);`

	checkIfOrdersTableExists      = `SELECT 'public.orders'::regclass;`
	
)

var message = []string{
	0: "DBStorage:unable to connect to database",
	1: "DBStorage:%v table has created",
	2: "DBStorage:unable to create %v table",
	3: "DBStorage:createOrUpdateIfExistsUsersTable error",
	4: "DBStorage:QueryRow failed: %v\n",
	5: "DBStorage:RowScan error",
	6: "DBStorage:time cannot be parsed",
	7: "DBStorage:createOrUpdateIfExistsWithdrawalsTable error",
}




func createTable(ctx context.Context, s DBStorage, checkTableSQL string, createTableSQL string) error {

	resp, err := s.pool.Exec(ctx, checkTableSQL)
	if err != nil {
		log.Println(message[2] + err.Error())
		//create Table
		resp, err = s.pool.Exec(ctx, createTableSQL)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message[1] + resp.String())
	} else {
		log.Println(message[2] + resp.String())
	}

	return err
}



//GetUser(ctx context.Context, name string) (u *schema.User, err error)
//SaveUser(ctx context.Context, u *schema.User) (err error)
//
//GetOrder(ctx context.Context, orderNumber int) (o *schema.Order, err error)
//SaveOrder(ctx context.Context, o schema.Order) (err error)
//GetOrdersList(ctx context.Context, userName string) (wl schema.Orders, err error)
//
//SaveWithdrawal(ctx context.Context, w schema.Withdrawal) (err error)
//GetWithdrawalsList(ctx context.Context, u schema.User) (wl *schema.Withdrawals, err error)

func (s DBStorage) GetUser(ctx context.Context, name string) (u *schema.User, err error) {
	if !s.connectDB(ctx) {
		return nil, errors.New(message[0])
	}
	defer s.conn.Release()
	d := dbUsers{userID: sql.NullString{String: name, Valid: true}}
	row := s.conn.QueryRow(ctx, selectLineUsersTable, &d.userID)
	err = row.Scan(&d.userID, &d.password, &d.accrual, &d.withdrawal)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, nil
		}
		return nil, err
	}
	return &schema.User{
		User:       d.userID.String,
		Password:   d.password.String,
		Accrual:    d.accrual.Float64,
		Withdrawal: d.withdrawal.Float64,
	}, nil
}

func (s DBStorage) SaveUser(ctx context.Context, u *schema.User) (err error) {
	if !s.connectDB(ctx) {
		return errors.New(message[0])
	}
	defer s.conn.Release()

	d := dbUsers{
		userID:    sql.NullString{String: u.User, Valid: true},
		password:   sql.NullString{String: u.Password, Valid: true},
		accrual:    sql.NullFloat64{Float64: u.Accrual, Valid: true},
		withdrawal: sql.NullFloat64{Float64: u.Withdrawal, Valid: true},
	}

	tag, err := s.conn.Exec(ctx, createOrUpdateIfExistsUsersTable, d.userID, d.password, d.accrual, d.withdrawal)
	logFatalf(message[3], err)
	log.Println(tag)
	return err
}

func (s DBStorage) GetOrder(ctx context.Context, orderNumber int64) (o *schema.Order, err error) {
	if !s.connectDB(ctx) {
		return nil, errors.New(message[0])
	}
	defer s.conn.Release()
	d := dbOrders{orderID: sql.NullInt64{Int64: orderNumber, Valid: true}}
	row := s.conn.QueryRow(ctx, selectLineOrdersTable, &d.orderID)
	err = row.Scan(&d.orderID, &d.userID, &d.status, &d.accrual, &d.createdAt)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return nil, err
	}
	created, err := time.Parse(time.RFC3339, d.createdAt.String)

	return &schema.Order{
		Order:   strconv.FormatInt(d.orderID.Int64, 10),
		User:    d.userID.String,
		Status:  schema.OrderStatus.ByCode[d.status.Int64].Text,
		Accrual: d.accrual.Float64,
		Created: schema.CreatedTime(created),
	}, nil
}
func (s DBStorage) SaveOrder(ctx context.Context, o schema.Order) (err error) {
	if !s.connectDB(ctx) {
		return errors.New(message[0])
	}

	orderInt, err := strconv.ParseInt(o.Order, 10, 64)
	if err != nil {
		log.Fatal(fmt.Errorf("error in converting order number %v to string:%w", o.Order, err))
	}

	d := &dbOrders{
		orderID:   sql.NullInt64{Int64: orderInt, Valid: true},
		userID:    sql.NullString{String: o.User, Valid: true},
		status:     sql.NullInt64{Int64: schema.OrderStatus.ByText[o.Status].Code, Valid: true},
		accrual:    sql.NullFloat64{Float64: o.Accrual, Valid: true},
		createdAt: sql.NullString{String: time.Time(o.Created).Format(time.RFC3339), Valid: true},
	}

	tag, err := s.conn.Exec(ctx, createOrUpdateIfExistsOrdersTable, d.orderID, d.userID, d.status, d.accrual, d.createdAt)
	logFatalf(message[3], err)
	log.Println(tag)
	return err
}

func (s DBStorage) GetOrdersList(ctx context.Context, userName string) (wl schema.Orders, err error) {
	if !s.connectDB(ctx) {
		return nil, errors.New(message[0])
	}
	defer s.conn.Release()

	wl = make(schema.Orders)

	d := dbOrders{userID: sql.NullString{String: userName, Valid: true}}

	rows, err := s.conn.Query(ctx, selectAllOrdersTableByUser, &d.userID)
	if err != nil {
		log.Printf(message[4], err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&d.orderID, &d.userID, &d.status, &d.accrual, &d.createdAt)
		logFatalf(message[5], err)
		created, err := time.Parse(time.RFC3339, d.createdAt.String)
		logFatalf(message[6], err)
		wl[d.orderID.Int64] = schema.Order{
			Order:   strconv.FormatInt(d.orderID.Int64, 10),
			User:    d.userID.String,
			Status:  schema.OrderStatus.ByCode[d.status.Int64].Text,
			Accrual: d.accrual.Float64,
			Created: schema.CreatedTime(created),
		}
	}

	return wl, nil
}

func (s DBStorage) GetNewOrdersList(ctx context.Context) (ol schema.Orders, err error) {
	if !s.connectDB(ctx) {
		return nil, errors.New(message[0])
	}
	defer s.conn.Release()

	ol = make(schema.Orders)

	d := dbOrders{status: sql.NullInt64{Int64: schema.OrderStatus.New.Code, Valid: true}}

	rows, err := s.conn.Query(ctx, selectAllOrdersTableByStatus, &d.status)
	if err != nil {
		log.Printf(message[4], err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&d.orderID, &d.userID, &d.status, &d.accrual, &d.createdAt)
		logFatalf(message[5], err)
		created, err := time.Parse(time.RFC3339, d.createdAt.String)
		logFatalf(message[6], err)
		ol[d.orderID.Int64] = schema.Order{
			Order:   strconv.FormatInt(d.orderID.Int64, 10),
			User:    d.userID.String,
			Status:  schema.OrderStatus.ByCode[d.status.Int64].Text,
			Accrual: d.accrual.Float64,
			Created: schema.CreatedTime(created),
		}
	}

	return ol, nil
}

func (s DBStorage) SaveWithdrawal(ctx context.Context, w schema.Withdrawal) (err error) {

	if !s.connectDB(ctx) {
		return errors.New(message[0])
	}
	defer s.conn.Release()

	d := dbWithdrawals{
		userID:    sql.NullString{String: w.User, Valid: true},
		createdAt: sql.NullString{String: time.Time(w.Processed).Format(time.RFC3339), Valid: true},
		orderID:   sql.NullString{String: w.Order, Valid: true},
		withdrawal: sql.NullFloat64{Float64: w.Withdrawal, Valid: true},
	}
	tag, err := s.conn.Exec(ctx, createOrUpdateIfExistsWithdrawalsTable, &d.userID, &d.createdAt, &d.orderID, &d.withdrawal)
	logFatalf(message[7], err)
	log.Println(tag)
	return err
}
func (s DBStorage) GetWithdrawalsList(ctx context.Context, username string) (wl *schema.Withdrawals, err error) {
	if !s.connectDB(ctx) {
		return nil, errors.New(message[0])
	}
	defer s.conn.Release()

	wl = new(schema.Withdrawals)

	d := dbWithdrawals{userID: sql.NullString{String: username, Valid: true}}

	rows, err := s.conn.Query(ctx, selectAllWithdrawalsTableByUser, &d.userID)
	if err != nil {
		log.Printf(message[4], err)
		return nil, err
	}
	log.Printf("getting withdrawals for user %v", d.userID)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&d.userID, &d.createdAt, &d.orderID, &d.withdrawal)
		logFatalf(message[5], err)
		created, err := time.Parse(time.RFC3339, d.createdAt.String)
		logFatalf(message[6], err)
		log.Printf("got withdrawal for user %v: %v", d.userID, d)

		w := schema.Withdrawal{
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
