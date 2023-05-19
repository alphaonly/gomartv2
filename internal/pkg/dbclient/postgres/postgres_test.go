package postgres

import (
	"testing"
)

func TestConnect(t *testing.T) {

	// t.Parallel()
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()

	// // given

	// mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	// // mockPool := pgxpoolmock.NewMockPgxIface(ctrl)

	// pool, err := pgxmock.NewPool()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer pool.Close()
	// var p pgxpool.Pool
	// dbClient := postgresClient{pool: mockPool}

	// open database stub
	//mock, err := pgxmock.NewPool()
	//if err != nil {
	//	t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	//}
	//
	//defer mock.Close()

}
