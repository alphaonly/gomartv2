package withdrawal

const (
	selectAllWithdrawalsTableByUser = `SELECT user_id,  uploaded_at,  order_id, withdrawal FROM public.withdrawals WHERE user_id = $1;`

	createOrUpdateIfExistsWithdrawalsTable = `
		INSERT INTO public.withdrawals (user_id, uploaded_at, order_id, withdrawal) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id,uploaded_at) DO UPDATE 
		  SET 	order_id   = $3,
		  		withdrawal = $4; 
		  `
	createWithdrawalsTable = `create table public.withdrawals
	(	user_id 		varchar(40) not null,
		uploaded_at 	TEXT 		not null,
		order_id   		varchar(40) not null,
		withdrawal 		double precision not null,
		primary key (user_id,uploaded_at)	
	);`

	checkIfWithdrawalsTableExists = `SELECT 'public.withdrawals'::regclass;`
)