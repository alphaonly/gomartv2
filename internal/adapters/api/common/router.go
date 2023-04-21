package api

import "github.com/go-chi/chi"

func NewRouter(h *handler) chi.Router {

	var (
		basicAuth = h.UserHandler.BasicAuth

		register       = h.Post(h.UserHandler.Register(nil))
		login          = h.Post(h.UserHandler.Login(nil))
		sendOrders     = h.Post(basicAuth(h.OrderHandler.PostOrders(nil)))
		withdraw       = h.Post(basicAuth(h.WithdrawalHandler.PostWithdraw(nil)))
		getOrders      = h.Get(basicAuth(h.OrderHandler.GetOrders(nil)))
		balance        = h.Get(basicAuth(h.OrderHandler.GetBalance(nil)))
		getWithdrawals = h.Get(basicAuth(h.WithdrawalHandler.GetWithdrawals(nil)))
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/check", h.HandleCheckHealth)

		r.Post("/api/user/register", register)
		r.Post("/api/user/login", login)
		r.Post("/api/user/orders", sendOrders)
		r.Post("/api/user/balance/withdraw", withdraw)

		r.Get("/api/user/orders", getOrders)
		r.Get("/api/user/balance", balance)
		r.Get("/api/user/withdrawals", getWithdrawals)

		//Mock for accrual system (in case similar addresses) returns +5 score
		r.Get("/api/orders/{number}", h.HandleGetOrderAccrual(nil))

	})

	return r
}
