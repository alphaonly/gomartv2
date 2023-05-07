package user

// BalanceResponseDTO - a structure for receiving user's balance
type BalanceResponseDTO struct {
	Current   float64 `json:"current"`   // current collected score
	Withdrawn float64 `json:"withdrawn"` // withdrawal score
}
