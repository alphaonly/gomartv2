package user

import "database/sql"

type DBUsersDTO struct {
	userID     sql.NullString
	password   sql.NullString
	accrual    sql.NullFloat64
	withdrawal sql.NullFloat64
}
