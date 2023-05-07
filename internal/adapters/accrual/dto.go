package accrual

// ResponseDTO - a transfer object structure for getting response from remote accrual service
type ResponseDTO struct {
	Order   int64   `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
