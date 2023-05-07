// Package withdrawal - this is a domain  part of withdrawal that contains withdrawal's model and service functionality.
package withdrawal

import (
	"time"

	"github.com/alphaonly/gomartv2/internal/schema"
)

type Withdrawal struct {
	User       string             `json:"user,omitempty"`
	Order      string             `json:"order"`
	Processed  schema.CreatedTime `json:"processed_at"`
	Withdrawal float64            `json:"withdrawal,omitempty"`
}
type Withdrawals []Withdrawal

func (ws Withdrawals) Response() (wrList *[]WithdrawalResponse) {
	wrList = new([]WithdrawalResponse)
	for _, w := range ws {
		wr := WithdrawalResponse{
			Order:      w.Order,
			Withdrawal: w.Withdrawal,
			Processed:  w.Processed,
		}
		*wrList = append(*wrList, wr)
	}
	return wrList
}

type ByTimeDescending Withdrawals

func (a ByTimeDescending) Len() int      { return len(a) }
func (a ByTimeDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTimeDescending) Less(i, j int) bool {
	return time.Time(a[i].Processed).Before(time.Time(a[i].Processed))
}
