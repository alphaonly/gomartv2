package user

const (
	selectLineUsersTable = `SELECT user_id, password, accrual, withdrawal FROM public.users WHERE user_id=$1;`

	createOrUpdateIfExistsUsersTable = `
	INSERT INTO public.users (user_id, password, accrual, withdrawal) 
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id) DO UPDATE 
  	SET password 	= $2,
	  	accrual 	= $3,
		withdrawal 	= $4; 
  	`
	createUsersTable = `create table public.users
	(	user_id varchar(40) not null primary key,
		password  TEXT not null,
		accrual double precision,
		withdrawal double precision 
	);`

	checkIfUsersTableExists = `SELECT 'public.users'::regclass;`
)

// -d=postgres://postgres:mypassword@localhost:5432/yandex
