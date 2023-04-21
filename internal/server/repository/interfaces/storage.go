package storage

// type Keeper interface {
// 	GetUser(ctx context.Context, name string) (u *user.User, err error)
// 	SaveUser(ctx context.Context, u *user.User) (err error)

// 	GetOrder(ctx context.Context, orderNumber int64) (o *order.Order, err error)
// 	SaveOrder(ctx context.Context, o order.Order) (err error)
// 	GetOrdersList(ctx context.Context, userName string) (ol order.Orders, err error)

// 	GetNewOrdersList(ctx context.Context) (ol order.Orders, err error)
// 	SaveWithdrawal(ctx context.Context, w withdrawal.Withdrawal) (err error)
// 	GetWithdrawalsList(ctx context.Context, userName string) (wl *withdrawal.Withdrawals, err error)
// }

// type Entity interface{
// }

// type User struct {
// 	Entity
// }

// type EntityKeeper interface {
// 	Get(ctx context.Context,id string) (u Entity,err error)
// }

// var e Entity = &User{}

// type UserRepository struct{

// }

// func (e *UserRepository)Get(ctx context.Context,id string) (user Entity,err error){

// return nil,nil

// }

// func test1(){
// 	var ek EntityKeeper= &UserRepository{}
// 	var u User

// 	e,err:= ek.Get(context.Background(),"")
// 	u=e.(User)
// 	if err!=nil{
// 		log.Fatal(u)
// 	}

// }
