package accrual

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/domain/user"
	"github.com/go-resty/resty/v2"
)

//Periodically checking orders' accrual from remote service

// Accrual - an interface for implementation accrual functionality
type Accrual interface {
	Run(ctx context.Context)
}

type accrual struct {
	serviceAddress string        // serviceAddress - address for request for a accrual
	requestTime    time.Duration //200 * time.Millisecond // requestTime - the time for repeat request
	OrderStorage   order.Storage
	UserStorage    user.Storage
}

// NewAccrual - a factory that bears new accrual entity
func NewAccrual(configuration *configuration.ServerConfiguration, orderStorage order.Storage, userStorage user.Storage) Accrual {
	return &accrual{
		serviceAddress: configuration.AccrualSystemAddress,
		requestTime:    time.Duration(configuration.AccrualTime) * time.Millisecond,
		OrderStorage:   orderStorage,
		UserStorage:    userStorage,
	}
}

// Run  - method to start requesting by an accrual entity
func (acr accrual) Run(ctx context.Context) {
	// ticker - it ticks once every acr.requestTime to repeat request
	ticker := time.NewTicker(acr.requestTime)
	baseURL, err := url.Parse(acr.serviceAddress)
	if err != nil {
		log.Fatal("unable to parse URL for accrual system")
	}

	httpc := resty.New().
		SetBaseURL(baseURL.String())

doItAGain:
	for {
		select {
		case <-ticker.C:
			// Getting New unprocessed orders to make a request to accrual system
			oList, err := acr.OrderStorage.GetNewOrdersList(ctx)
			if err != nil {
				log.Fatal("can not get new orders list")
			}

			for orderNumber, orderData := range oList {

				orderNumberStr := strconv.Itoa(int(orderNumber))
				req := httpc.R().
					SetHeader("Accept", "application/json")

				response := order.AccrualResponse{}
				resp, err := req.
					SetResult(&response).
					Get("api/orders/" + orderNumberStr)
				if err != nil {
					log.Printf("order %v response error: %v", orderNumber, err)
					continue
				}
				log.Printf("order %v response from accrual: %v", orderNumber, resp)
				if response.Status != order.ProcessedOrder.Text {
					log.Printf("order %v response status type %v, continue", orderNumber, resp.Status())
					continue
				}

				orderData.Accrual = response.Accrual
				orderData.Status = order.ProcessedOrder.Text
				log.Printf("Saving processed order:%v", orderData)

				err = acr.OrderStorage.SaveOrder(ctx, orderData)
				if err != nil {
					log.Fatal("unable to save order")
				}
				log.Printf("Processed order saved from accrual:%v", orderData)
				log.Printf("Update user balance with processed order:%v", orderNumber)
				//Update balance in case of order accrual greater than zero
				if orderData.Accrual > 0 {

					u, err := acr.UserStorage.GetUser(ctx, orderData.User)
					if err != nil {
						log.Fatalf("Error in getting user %v data: %v", orderData.User, err.Error())
					}
					if u == nil {
						log.Fatalf("Data inconsistency with there is no user %v, but there is order %v with the user", orderData.User, orderNumber)
					}
					u.Accrual += orderData.Accrual
					err = acr.UserStorage.SaveUser(ctx, u)
					if err != nil {
						log.Fatalf("Unable to save user %v with updated accrual %v: %v", u.User, u.Accrual, err.Error())
					}
					log.Printf("Updated user:%v", u)

				}
			}

		case <-ctx.Done():
			break doItAGain
		}
	}

}

func (acr accrual) sendRequest(ctx context.Context) {}

func (acr accrual) GetResponse(ctx context.Context) {}
