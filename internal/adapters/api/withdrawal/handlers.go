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
	PostWithdraw() http.HandlerFunc
	GetWithdrawals() http.HandlerFunc
}

type handler struct {
	Storage       withdrawal.Storage
	Service       withdrawal.Service
	OrderService  order.Service
	Configuration *configuration.ServerConfiguration
}

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

func (h *handler) PostWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostUserBalanceWithdraw invoked")
		//Get parameters from previous handler
		userName, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HttpError(w, fmt.Errorf("cannot get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		requestByteData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unrecognized json request ", http.StatusBadRequest)
			return
		}
		userWithdrawalRequest := withdrawal.UserWithdrawalRequestDTO{}
		err = json.Unmarshal(requestByteData, &userWithdrawalRequest)
		if err != nil {
			http.Error(w, "Error json-marshal request data", http.StatusBadRequest)
			return
		}
		err = h.Service.MakeUserWithdrawal(r.Context(), string(userName), userWithdrawalRequest)
		if err != nil {
			if errors.Is(err, withdrawal.ErrNoFunds) {
				api.HttpErrorW(w, "make withdrawal error", err, http.StatusPaymentRequired)
				return
			}
			if errors.Is(err, withdrawal.ErrOrderInvalid) {
				api.HttpErrorW(w, "order number invalid", err, http.StatusUnprocessableEntity)
				return
			}
			if errors.Is(err, withdrawal.ErrInternal) {
				api.HttpErrorW(w, "internal error", err, http.StatusInternalServerError)
				return
			}
		}
		//Response
		w.WriteHeader(http.StatusOK)
	}
}
func (h *handler) GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetUserWithdrawals invoked")
		//Get parameters from previous handler
		userName, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HttpError(w, fmt.Errorf("can not get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		wList, err := h.Service.GetUsersWithdrawals(r.Context(), string(userName))
		if err != nil {
			if errors.Is(err, withdrawal.ErrInternal) {
				api.HttpErrorW(w, "internal error", err, http.StatusInternalServerError)
				return
			}
			if errors.Is(err, withdrawal.ErrNoWithdrawal) {
				api.HttpErrorW(w, "no withdrawals", err, http.StatusNoContent)
				return
			}
		}
		response := wList.Response()
		log.Printf("return withdrawals response list: %v", response)
		//Response
		bytes, err := json.Marshal(response)
		if err != nil {
			api.HttpErrorW(w, fmt.Sprintf("user %v withdrawals list json marshal error", userName), err, http.StatusInternalServerError)
			return
		}
		log.Printf("return withdrawals list in JSON: %v", string(bytes))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(bytes)
		if err != nil {
			api.HttpErrorW(w, fmt.Sprintf("user %v withdrawals list write response error", userName), err, http.StatusInternalServerError)
			return
		}
	}
}
