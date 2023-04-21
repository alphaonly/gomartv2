package order

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/alphaonly/gomartv2/internal/schema"
)

type Order struct {
	Order   string             `json:"number"`
	User    string             `json:"user"`
	Status  string             `json:"status,omitempty"`
	Accrual float64            `json:"accrual,omitempty"`
	Created schema.CreatedTime `json:"uploaded_at"`
}

type orderType struct {
	Code int64
	Text string
}

var (
	NewOrder        = orderType{1, "NEW"}
	ProcessingOrder = orderType{2, "PROCESSING"}
	InvalidOrder    = orderType{3, "INVALID"}
	ProcessedOrder  = orderType{4, "PROCESSED"}
)

type Orders map[int64]Order

func (o Orders) MarshalJSON() ([]byte, error) {

	oArray := make([]Order, len(o))
	i := 0
	for k, v := range o {
		oNumb := strconv.FormatInt(k, 10)
		oArray[i] = Order{
			Order:   oNumb,
			Status:  v.Status,
			Accrual: v.Accrual,
			Created: v.Created,
		}
		i++
	}
	bytes, err := json.Marshal(&oArray)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
func (o Orders) UnmarshalJSON(b []byte) error {
	var oArray []Order
	if err := json.Unmarshal(b, &oArray); err != nil {
		return err
	}
	for _, v := range oArray {
		OrderInt, err := strconv.ParseInt(v.Order, 10, 64)
		if err != nil {
			log.Fatal(fmt.Errorf("cannot convert order number %v to string: %w", OrderInt, err))
		}
		o[OrderInt] = v
	}
	return nil
}
