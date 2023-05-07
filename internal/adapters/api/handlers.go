package api

import (
	"encoding/json"
	"fmt"
	"github.com/alphaonly/gomartv2/internal/pkg/common/logging"
	"log"
	"net/http"
	"strconv"

	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/go-chi/chi"
)

// Handler - an interface that implements common app handlers.
type Handler interface {
	Health() http.HandlerFunc                        // a function to implement HTTP GET request to check server alive
	Post(next http.Handler) http.HandlerFunc         // a technical function handler to implement a check whether the request is  a POST request
	Get(next http.Handler) http.HandlerFunc          // a technical function handler to implement a check whether the request  is a GET request
	AccrualScore(next http.Handler) http.HandlerFunc // a technical function handler to implement an accrual system mock for quick test
}

// NewHandler - it is a factory that returns an instance of common Handler implementation.
func NewHandler(
	configuration *configuration.ServerConfiguration) Handler {

	return &handler{
		Configuration: configuration}
}

type handler struct {
	Configuration *configuration.ServerConfiguration // a pointer to a server configuration
}

// Health - a function to implement HTTP GET request to check server alive
func (h *handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("200 OK"))
		logging.LogFatal(err)
	}
}

// Get - a technical function handler to implement a check whether the request is  a POST request
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

// Post - a technical function handler to implement a check whether the request is  a POST request
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

// AccrualScore - a technical function handler to implement an accrual system mock for quick test
func (h *handler) AccrualScore(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetOrderAccrual invoked")
		//Handling
		orderNumberStr := chi.URLParam(r, "number")
		if orderNumberStr == "" {
			HTTPError(w, fmt.Errorf("order number  %v is empty", orderNumberStr), http.StatusBadRequest)
			return
		}

		_, err := strconv.ParseInt(orderNumberStr, 10, 64)
		if err != nil {
			HTTPError(w, fmt.Errorf("order number  %v is bad format", orderNumberStr), http.StatusBadRequest)
			return
		}

		accrual := 5.3

		OrderAccrualResponse := order.AccrualResponse{
			Order:   orderNumberStr,
			Status:  "PROCESSED",
			Accrual: accrual,
		}

		//Response
		bytes, err := json.Marshal(&OrderAccrualResponse)
		if err != nil {
			HTTPErrorW(w, "order accrual response json marshal error", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "apilication/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(bytes)
		if err != nil {
			HTTPErrorW(w, "order accrual response write response error", err, http.StatusInternalServerError)
			return
		}

	}
}
