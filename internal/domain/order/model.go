// Package order - this is a domain  part of order that contains orders' model and service functionality.
package order

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/alphaonly/gomartv2/internal/schema"
)

// Order - a domain structure for order entity
type Order struct {
	Order   string             `json:"number"`            // order's id
	User    string             `json:"user"`              // user-owner of order
	Status  string             `json:"status,omitempty"`  //status of order processing
	Accrual float64            `json:"accrual,omitempty"` // accrual assigned by remote accrual system
	Created schema.CreatedTime `json:"uploaded_at"`       // date and time of order's creation
}

type orderType struct {
	Code int64
	Text string
}

// Constants of order's processing status variants
var (
	NewOrder        = orderType{1, "NEW"}
	ProcessingOrder = orderType{2, "PROCESSING"}
	InvalidOrder    = orderType{3, "INVALID"}
	ProcessedOrder  = orderType{4, "PROCESSED"}
)

// OrderTypesByCode - Constants of order's processing status values given in int codes
var OrderTypesByCode = map[int64]orderType{
	NewOrder.Code:        NewOrder,
	ProcessingOrder.Code: ProcessingOrder,
	InvalidOrder.Code:    InvalidOrder,
	ProcessedOrder.Code:  ProcessedOrder}

// OrderTypesByText - Constants of order's processing status values given in int strings
var OrderTypesByText = map[string]orderType{
	NewOrder.Text:        NewOrder,
	ProcessingOrder.Text: ProcessingOrder,
	InvalidOrder.Text:    InvalidOrder,
	ProcessedOrder.Text:  ProcessedOrder}

// Orders - a type that describes a hashmap for Order type
type Orders map[int64]Order

// MarshalJSON - a function to encode a list of orders to JSON format
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

// UnmarshalJSON - a function to decode a list of orders from JSON format
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
