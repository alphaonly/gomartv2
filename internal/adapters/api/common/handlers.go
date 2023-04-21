package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alphaonly/gomartv2/internal/adapters/api/app"
	orderApi "github.com/alphaonly/gomartv2/internal/adapters/api/order"
	userApi "github.com/alphaonly/gomartv2/internal/adapters/api/user"
	withdrawalApi "github.com/alphaonly/gomartv2/internal/adapters/api/withdrawal"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/go-chi/chi"
)

type Handler interface {
	HandleCheckHealth(w http.ResponseWriter, r *http.Request)
	BasicUserAuthorization(next http.Handler) http.HandlerFunc
	POST(next http.Handler) http.HandlerFunc
	GET(next http.Handler) http.HandlerFunc
}

type handler struct {
	UserHandler       userApi.Handler
	OrderHandler      orderApi.Handler
	WithdrawalHandler withdrawalApi.Handler
	Configuration     configuration.ServerConfiguration
}

func (h *handler) HandleCheckHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		w.WriteHeader(http.StatusOK)

	}
}

func (h *handler) Get(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetValidation invoked")
		//Validation
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
			return
		}
		if next != nil {
			//call further handler with context parameters
			next.ServeHTTP(w, r)
			return
		}
	}
}
func (h *handler) Post(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostValidation invoked")
		//Validation
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
			return
		}
		if next != nil {
			//call further handler with context parameters
			next.ServeHTTP(w, r)
			return
		}
	}
}

func (h *handler) HandleGetOrderAccrual(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetOrderAccrual invoked")
		//Handling
		orderNumberStr := chi.URLParam(r, "number")
		if orderNumberStr == "" {
			app.HttpError(w, fmt.Errorf("order number  %v is empty", orderNumberStr), http.StatusBadRequest)
			return
		}

		_, err := strconv.ParseInt(orderNumberStr, 10, 64)
		if err != nil {
			app.HttpError(w, fmt.Errorf("order number  %v is bad format", orderNumberStr), http.StatusBadRequest)
			return
		}

		accrual := 5.3

		OrderAccrualResponse := order.OrderAccrualResponse{
			Order:   orderNumberStr,
			Status:  "PROCESSED",
			Accrual: accrual,
		}

		//Response
		bytes, err := json.Marshal(&OrderAccrualResponse)
		if err != nil {
			app.HttpErrorW(w, "order accrual response json marshal error", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(bytes)
		if err != nil {
			app.HttpErrorW(w, "order accrual response write response error", err, http.StatusInternalServerError)
			return
		}

	}
}
