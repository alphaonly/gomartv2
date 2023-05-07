// Package order - this database part of order entity that contains orders' functions with a database communication
package order

import "database/sql"

// DBOrdersDTO - a transfer object structure to communicate order's data to postgres library pgx
type DBOrdersDTO struct {
	orderID   sql.NullInt64   // order ID
	userID    sql.NullString  // user ID
	status    sql.NullInt64   // order's processing status
	accrual   sql.NullFloat64 // accrual for the order
	createdAt sql.NullString  // date and time of order creation
}
