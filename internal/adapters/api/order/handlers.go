package order

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alphaonly/gomartv2/internal/adapters/api"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/order"
	"github.com/alphaonly/gomartv2/internal/domain/user"
	"github.com/alphaonly/gomartv2/internal/schema"
)

type Handler interface {
	PostOrders(next http.Handler) http.HandlerFunc
	GetOrders(next http.Handler) http.HandlerFunc
	GetBalance(next http.Handler) http.HandlerFunc
}
type handler struct {
	Storage       order.Storage
	Service       order.Service
	UserService   user.Service
	Configuration *configuration.ServerConfiguration
}

func NewHandler(storage order.Storage, service order.Service, configuration *configuration.ServerConfiguration) Handler {
	return &handler{
		Storage:       storage,
		Service:       service,
		Configuration: configuration,
	}
}

func (h *handler) PostOrders(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostUserOrders invoked")
		//Get parameters from previous handler
		user, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HttpError(w, fmt.Errorf("cannot get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		OrderNumberByte, err := io.ReadAll(r.Body)
		if err != nil {
			api.HttpError(w, fmt.Errorf("unrecognized body body %w", err), http.StatusBadRequest)
			return
		}

		orderNumber, err := h.Service.ValidateOrderNumber(r.Context(), string(OrderNumberByte), string(user))
		if err != nil {
			if strings.Contains(err.Error(), "400") {
				api.HttpErrorW(w, fmt.Sprintf("order number  %v insufficient format", orderNumber), err, http.StatusBadRequest)
				return
			}
			if strings.Contains(err.Error(), "422") {
				api.HttpErrorW(w, fmt.Sprintf("order %v insufficient format", orderNumber), err, http.StatusUnprocessableEntity)
				return
			}
			if strings.Contains(err.Error(), "409") {
				api.HttpErrorW(w, fmt.Sprintf("order %v exists", orderNumber), err, http.StatusConflict)
				return
			}
			if strings.Contains(err.Error(), "200") {
				log.Printf("order %v exists: %v", orderNumber, err.Error())
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		//Create object for a new order
		o := order.Order{
			Order:   string(OrderNumberByte),
			User:    string(user),
			Status:  order.NewOrder.Text,
			Created: schema.CreatedTime(time.Now()),
		}
		err = h.Storage.SaveOrder(r.Context(), o)
		if err != nil {
			api.HttpErrorW(w, fmt.Sprintf("order's number %v not saved", orderNumber), err, http.StatusInternalServerError)
			return
		}
		//Response
		w.WriteHeader(http.StatusAccepted)
	}
}
func (h *handler) GetOrders(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetUserOrders invoked")

		//Get parameters from previous handler
		userName, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HttpError(w, fmt.Errorf("cannot get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		orderList, err := h.Service.GetUsersOrders(r.Context(), string(userName))
		if err != nil {
			if strings.Contains(err.Error(), "204") {
				api.HttpErrorW(w, fmt.Sprintf("No orders for user %v", userName), err, http.StatusNoContent)
				return
			}
		}
		//Response
		bytes, err := json.Marshal(orderList)
		if err != nil {
			api.HttpErrorW(w, fmt.Sprintf("user %v order list json marshal error", userName), err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(bytes)
		if err != nil {
			api.HttpErrorW(w, fmt.Sprintf("user %v HandleGetUserOrders write response error", userName), err, http.StatusInternalServerError)
			return
		}
	}
}
func (h *handler) GetBalance(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleGetUserBalance invoked")
		//Get parameters from previous handler
		userName, err := api.GetPreviousParameter[schema.CtxUName, schema.ContextKey](r, schema.CtxKeyUName)
		if err != nil {
			api.HttpError(w, fmt.Errorf("cannot get userName from context %w", err), http.StatusInternalServerError)
			return
		}
		//Handling
		balance, err := h.UserService.GetUserBalance(r.Context(), string(userName))
		if err != nil {
			api.HttpError(w, fmt.Errorf("cannot get user data by userName %v from context %w", userName, err), http.StatusInternalServerError)
			return
		}
		log.Printf("Got balance %v for user %v ", balance, userName)
		//Response
		bytes, err := json.Marshal(balance)
		if err != nil {
			api.HttpErrorW(w, fmt.Sprintf("user %v balance json marshal error", userName), err, http.StatusInternalServerError)
			return
		}
		log.Printf("Write response balance json:%v ", string(bytes))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write(bytes)
		if err != nil {
			api.HttpErrorW(w, fmt.Sprintf("user %v balance write response error", userName), err, http.StatusInternalServerError)
			return
		}
	}
}
