// Package withdrawal - this API part of withdrawal entity that contains withdrawal' handlers functionality.
package withdrawal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alphaonly/gomartv2/internal/adapters/api"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/domain/withdrawal"
	"github.com/alphaonly/gomartv2/internal/schema"
	"io"
	"log"
	"net/http"
)

type Handler interface {
	PostWithdraw() http.HandlerFunc   // a function to implement HTTP POST withdrawal data handler to make a withdrawal
	GetWithdrawals() http.HandlerFunc // a function to implement HTTP GET withdrawal data handler
}

type handler struct {
	Storage       withdrawal.Storage                 // a pointer to withdrawal storage interface implementation
	Service       withdrawal.Service                 // a pointer to withdrawal service interface implementation
	OrderService  order.Service                      // a pointer to order service interface implementation
	Configuration *configuration.ServerConfiguration // a pointer to a server configuration
}

// NewHandler - it is a factory that returns an instance of withdrawal's Handler implementation.
func NewHandler(
	storage withdrawal.Storage,
	service withdrawal.Service,
	orderService order.Service,
	configuration *configuration.ServerConfiguration) Handler {
	return &handler{
		Storage:       storage,
		Service:       service,
		OrderService:  orderService,
		Configuration: configuration,
	}
}

// PostWithdraw - a handler function to implement HTTP POST request to handle making a withdrawal
func (h *handler) PostWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostUserBalanceWithdraw invoked")
		//Get parameters from previous handler
		userName, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HTTPError(w, fmt.Errorf("cannot get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		requestByteData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unrecognized json request ", http.StatusBadRequest)
			return
		}
		userWithdrawalRequest := withdrawal.WithdrawalRequestDTO{}
		err = json.Unmarshal(requestByteData, &userWithdrawalRequest)
		if err != nil {
			http.Error(w, "Error json-marshal request data", http.StatusBadRequest)
			return
		}
		err = h.Service.MakeUserWithdrawal(r.Context(), string(userName), userWithdrawalRequest)
		if err != nil {
			if errors.Is(err, withdrawal.ErrNoFunds) {
				api.HTTPErrorW(w, "make withdrawal error", err, http.StatusPaymentRequired)
				return
			}
			if errors.Is(err, withdrawal.ErrOrderInvalid) {
				api.HTTPErrorW(w, "order number invalid", err, http.StatusUnprocessableEntity)
				return
			}
			if errors.Is(err, withdrawal.ErrInternal) {
				api.HTTPErrorW(w, "internal error", err, http.StatusInternalServerError)
				return
			}
		}
		//Response
		w.WriteHeader(http.StatusOK)
	}
}

// GetWithdrawals - a handler function to implement HTTP GET request to handle getting withdrawal data
func (h *handler) GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetUserWithdrawals invoked")
		//Get parameters from previous handler
		userName, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HTTPError(w, fmt.Errorf("can not get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		wList, err := h.Service.GetUsersWithdrawals(r.Context(), string(userName))
		if err != nil {
			if errors.Is(err, withdrawal.ErrInternal) {
				api.HTTPErrorW(w, "internal error", err, http.StatusInternalServerError)
				return
			}
			if errors.Is(err, withdrawal.ErrNoWithdrawal) {
				api.HTTPErrorW(w, "no withdrawals", err, http.StatusNoContent)
				return
			}
		}
		response := wList.Response()
		log.Printf("return withdrawals response list: %v", response)
		//Response
		bytes, err := json.Marshal(response)
		if err != nil {
			api.HTTPErrorW(w, fmt.Sprintf("user %v withdrawals list json marshal error", userName), err, http.StatusInternalServerError)
			return
		}
		log.Printf("return withdrawals list in JSON: %v", string(bytes))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(bytes)
		if err != nil {
			api.HTTPErrorW(w, fmt.Sprintf("user %v withdrawals list write response error", userName), err, http.StatusInternalServerError)
			return
		}
	}
}
