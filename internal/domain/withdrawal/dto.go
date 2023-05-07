package withdrawal

import "github.com/alphaonly/gomartv2/internal/schema"

type WithdrawalRequestDTO struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawalResponse struct {
	Order      string             `json:"order"`
	Withdrawal float64            `json:"sum"`
	Processed  schema.CreatedTime `json:"processed_at"`
}
