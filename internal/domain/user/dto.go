package user

type BalanceResponseDTO struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
