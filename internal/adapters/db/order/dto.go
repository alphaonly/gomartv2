// Package order - this database part of order entity that contains orders' functions with a database communication
package order

import "database/sql"

// DBOrdersDTO - a transfer object structure to communicate order's data to postgres library pgx
type DBOrdersDTO struct {
	orderID   sql.NullInt64
	userID    sql.NullString
	status    sql.NullInt64
	accrual   sql.NullFloat64
	createdAt sql.NullString
}
