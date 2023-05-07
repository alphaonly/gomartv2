// Package withdrawal - this database part of withdrawal entity that contains orders' functions with a database communication
package withdrawal

import "database/sql"

type DBWithdrawalsDTO struct {
	userID     sql.NullString
	createdAt  sql.NullString
	orderID    sql.NullString
	withdrawal sql.NullFloat64
}
