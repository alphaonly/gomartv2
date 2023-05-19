package postgres

import (
	"context"
	"log"
)

var Message = []string{
	0:  "postgres client:unable to connect to database",
	1:  "postgres client:%v table has created",
	2:  "postgres client:unable to create %v table",
	3:  "postgres client:createOrUpdateIfExistsUsersTable error",
	4:  "postgres client:QueryRow failed: %v\n",
	5:  "postgres client:RowScan error",
	6:  "postgres client:time cannot be parsed",
	7:  "postgres client:createOrUpdateIfExistsWithdrawalsTable error",
	8:  "postgres client:unable to get postgres conn",
	9:  "postgres client:unable to get postgres conn pull",
	10: "postgres client:unable to create or check tables",
}

const (
	selectLineUsersTable            = `SELECT user_id, password, accrual, withdrawal FROM public.users WHERE user_id=$1;`
	selectLineOrdersTable           = `SELECT order_id, user_id, status, accrual, uploaded_at FROM public.orders WHERE order_id=$1;`
	selectAllOrdersTableByUser      = `SELECT order_id, user_id, status, accrual, uploaded_at FROM public.orders WHERE user_id = $1;`
	selectAllOrdersTableByStatus    = `SELECT order_id, user_id, status, accrual, uploaded_at  FROM public.orders WHERE status = $1;`
	selectAllWithdrawalsTableByUser = `SELECT user_id,  uploaded_at,  order_id, withdrawal FROM public.withdrawals WHERE user_id = $1;`

	createOrUpdateIfExistsUsersTable = `
	INSERT INTO public.users (user_id, password, accrual, withdrawal) 
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id) DO UPDATE 
  	SET password 	= $2,
	  	accrual 	= $3,
		withdrawal 	= $4; 
  	`
	createOrUpdateIfExistsOrdersTable = `
	  INSERT INTO public.orders (order_id, user_id, status,accrual,uploaded_at) 
	  VALUES ($1, $2, $3,$4, $5)
	  ON CONFLICT (order_id,user_id) DO UPDATE 
		SET status 		= $3,
		    accrual 	= $4,
			uploaded_at = $5; 
		`
	createOrUpdateIfExistsWithdrawalsTable = `
		INSERT INTO public.withdrawals (user_id, uploaded_at, order_id, withdrawal) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id,uploaded_at) DO UPDATE 
		  SET 	order_id   = $3,
		  		withdrawal = $4; 
		  `
	createUsersTable = `create table public.users
	(	user_id varchar(40) not null primary key,
		password  TEXT not null,
		accrual double precision,
		withdrawal double precision 
	);`
	createOrdersTable = `create table public.orders
	(	order_id bigint not null, 
		user_id varchar(40) not null,
		status integer,		
		accrual double precision,
		uploaded_at TEXT not null, 
		primary key (order_id,user_id)
	);`
	createWithdrawalsTable = `create table public.withdrawals
	(	user_id 		varchar(40) not null,
		uploaded_at 	TEXT 		not null,
		order_id   		varchar(40) not null,
		withdrawal 		double precision not null,
		primary key (user_id,uploaded_at)	
	);`

	checkIfUsersTableExists       = `SELECT 'public.users'::regclass;`
	checkIfOrdersTableExists      = `SELECT 'public.orders'::regclass;`
	checkIfWithdrawalsTableExists = `SELECT 'public.withdrawals'::regclass;`
)

func CreateTable(ctx context.Context, s postgresClient, checkTableSQL string, createTableSQL string) error {

	resp, err := s.pool.Exec(ctx, checkTableSQL)
	if err != nil {
		log.Println(Message[2] + err.Error())
		//create Table
		resp, err = s.pool.Exec(ctx, createTableSQL)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(Message[1] + resp.String())
	} else {
		log.Println(Message[2] + resp.String())
	}

	return err
}
