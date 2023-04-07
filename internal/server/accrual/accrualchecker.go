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

	baseURL := url.URL{
		Scheme: "http",
		Host:   c.serviceAddress,
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

			for orderNumber, data := range oList {

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
				if response.Status != "PROCESSED" {
					continue
				}

				data.Accrual = response.Accrual
				data.Status = schema.OrderStatus["PROCESSED"]

				err = c.storage.SaveOrder(ctx, data)
				if err != nil {
					log.Fatal("unable to save order")
				}
			}

		case <-ctx.Done():
			break doItAGain
		}
	}

}
