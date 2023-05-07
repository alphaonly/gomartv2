package order

// Constants that are SELECT REQUESTS to POSTGRES database with given parameters
const (
	selectLineOrdersTable        = `SELECT order_id, user_id, status, accrual, uploaded_at FROM public.orders WHERE order_id=$1;`
	selectAllOrdersTableByUser   = `SELECT order_id, user_id, status, accrual, uploaded_at FROM public.orders WHERE user_id = $1;`
	selectAllOrdersTableByStatus = `SELECT order_id, user_id, status, accrual, uploaded_at FROM public.orders WHERE status = $1;`

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

	checkIfOrdersTableExists = `SELECT 'public.orders'::regclass;`
)

// -d=postgres://postgres:mypassword@localhost:5432/yandex
