package withdrawal

import "database/sql"

type DBWithdrawalsDTO struct {
	userID     sql.NullString
	createdAt  sql.NullString
	orderID    sql.NullString
	withdrawal sql.NullFloat64
}
