package accrual

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/alphaonly/gomartv2/internal/schema"
	storage "github.com/alphaonly/gomartv2/internal/server/storage/interfaces"
	"github.com/go-resty/resty/v2"
)

//Periodically checking orders' accrual from remote service

type Configuration struct {
}

func NewChecker(serviceAddress string, requestTime int64, storage storage.Storage) (c *Checker) {
	return &Checker{
		serviceAddress: serviceAddress,
		requestTime:    time.Duration(requestTime) * time.Millisecond,
		storage:        storage,
	}
}

type Checker struct {
	serviceAddress string
	requestTime    time.Duration //200 * time.Millisecond
	storage        storage.Storage
}
type Response struct {
	Order   int64   `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (c Checker) Run(ctx context.Context) {
	ticker := time.NewTicker(c.requestTime)
	baseURL, err := url.Parse(c.serviceAddress)
	if err != nil {
		log.Fatal("unable to parse URL for accrual system")
	}

	httpc := resty.New().
		SetBaseURL(baseURL.String())

doItAGain:
	for {
		select {
		case <-ticker.C:
			//Getting New unprocessed orders to make a request to accrual system
			oList, err := c.storage.GetNewOrdersList(ctx)
			if err != nil {
				log.Fatal("can not get new orders list")
			}

			for orderNumber, order := range oList {

				orderNumberStr := strconv.Itoa(int(orderNumber))
				req := httpc.R().
					SetHeader("Accept", "application/json")

				response := schema.OrderAccrualResponse{}
				resp, err := req.
					SetResult(&response).
					Get("api/orders/" + orderNumberStr)
				if err != nil {
					log.Printf("order %v response error: %v", orderNumber, err)
					continue
				}
				log.Printf("order %v response from accrual: %v", orderNumber, resp)
				if response.Status != schema.OrderStatus.ByText["PROCESSED"].Text {
					log.Printf("order %v response status type %v, continue", orderNumber, resp.Status())
					continue
				}

				order.Accrual = response.Accrual
				order.Status = schema.OrderStatus.ByText["PROCESSED"].Text
				log.Printf("Saving processed order:%v", order)

				err = c.storage.SaveOrder(ctx, order)
				if err != nil {
					log.Fatal("unable to save order")
				}
				log.Printf("Processed order saved from accrual:%v", order)
				log.Printf("Update user balance with processed order:%v", orderNumber)
				//Update balance in case of order accrual greater than zero
				if order.Accrual > 0 {

					u, err := c.storage.GetUser(ctx, order.User)
					if err != nil {
						log.Fatalf("Error in getting user %v data: %v", order.User, err.Error())
					}
					if u == nil {
						log.Fatalf("Data inconsistency with there is no user %v, but there is order %v with the user", order.User, orderNumber)
					}
					u.Accrual += order.Accrual
					err = c.storage.SaveUser(ctx, u)
					if err != nil {
						log.Fatalf("Unable to save user %v with updated accrual %v: %v", u.User, u.Accrual,err.Error())
					}
					log.Printf("Updated user:%v", u)

				}
			}

		case <-ctx.Done():
			break doItAGain
		}
	}

}
