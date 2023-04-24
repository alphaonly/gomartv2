package order

import "database/sql"

type DBOrdersDTO struct {
	orderID   sql.NullInt64
	userID    sql.NullString
	status    sql.NullInt64
	accrual   sql.NullFloat64
	createdAt sql.NullString
}
