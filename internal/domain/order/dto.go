package order

// AccrualResponse - a struct for receiving data from a remote accrual score system
type AccrualResponse struct {
	Order   string  `json:"order"`   // order id
	Status  string  `json:"status"`  // order status
	Accrual float64 `json:"accrual"` // accrual score for order
}
