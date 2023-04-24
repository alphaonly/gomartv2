package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/alphaonly/gomartv2/internal/adapters/api"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/domain/user"
	"github.com/alphaonly/gomartv2/internal/schema"
)

type Handler interface {
	Register(next http.Handler) http.HandlerFunc
	Login(next http.Handler) http.HandlerFunc
	BasicAuth(next http.Handler) http.HandlerFunc
}

type handler struct {
	Storage       user.Storage
	Service       user.Service
	Configuration *configuration.ServerConfiguration
}

func NewHandler(storage user.Storage, service user.Service, configuration *configuration.ServerConfiguration) Handler {
	return &handler{
		Storage:       storage,
		Service:       service,
		Configuration: configuration,
	}
}

func (h *handler) Register(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostUserRegister invoked")

		//Handling body
		requestByteData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unrecognized json request ", http.StatusBadRequest)
			return
		}
		u := new(user.User)
		err = json.Unmarshal(requestByteData, u)
		if err != nil {
			http.Error(w, "Error json-marshal request data", http.StatusBadRequest)
			return
		}
		//Logic
		err = h.Service.RegisterUser(r.Context(), u)
		if err != nil {
			if strings.Contains(err.Error(), "400") {
				http.Error(w, "login "+u.User+": bad request", http.StatusBadRequest)
				return
			}
			if strings.Contains(err.Error(), "409") {
				http.Error(w, "login "+u.User+"is occupied", http.StatusConflict)
				return
			}
			http.Error(w, "login "+u.User+"register internal error", http.StatusInternalServerError)
			return
		}
		//Response
		log.Printf("Respond in header basic authorization: user:%v password: %v", u.User, u.Password)
		w.Header().Add("Authorization", "Basic "+api.BasicAuth(u.User, u.Password))
		w.WriteHeader(http.StatusOK)

	}
}

func (h *handler) Login(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandlePostUserLogin invoked")

		//Handling body
		requestByteData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unrecognized json request ", http.StatusBadRequest)
			return
		}
		u := new(user.User)
		err = json.Unmarshal(requestByteData, u)
		if err != nil {
			http.Error(w, "error json-marshal request data", http.StatusBadRequest)
			return
		}
		//Logic
		err = h.Service.AuthenticateUser(r.Context(), u)
		if err != nil {
			if strings.Contains(err.Error(), "400") {
				http.Error(w, "login "+u.User+": bad request", http.StatusBadRequest)
				return
			}
			if strings.Contains(err.Error(), "401") {
				api.HttpErrorW(w, "authorization error", err, http.StatusUnauthorized)
				return
			}
			if strings.Contains(err.Error(), "409") {
				http.Error(w, "login "+u.User+"is occupied", http.StatusConflict)
				return
			}
			api.HttpErrorW(w, "login "+u.User+"register internal error", err, http.StatusInternalServerError)
			return
		}
		//Response
		log.Printf("Respond in header basic authorization: user:%v password: %v", u.User, u.Password)
		w.Header().Add("Authorization", "Basic "+api.BasicAuth(u.User, u.Password))
		w.WriteHeader(http.StatusOK)
	}
}

func (h *handler) BasicAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("BasicUserAuthorization invoked")
		//Basic authentication
		userBA, passBA, ok := r.BasicAuth()
		if !ok {
			api.HttpError(w, fmt.Errorf("basic authentication is not ok"), http.StatusUnauthorized)
			return
		}
		log.Printf("basic authorization check: user: %v, password: %v", userBA, passBA)

		var err error
		ok, err = h.Service.CheckIfUserAuthorized(r.Context(), userBA, passBA)
		if err != nil {
			if strings.Contains(err.Error(), "400") {
				api.HttpError(w, fmt.Errorf("login %v: bad request %w", userBA, err), http.StatusBadRequest)
				return
			}
			if strings.Contains(err.Error(), "500") {
				api.HttpError(w, fmt.Errorf("login %v: server internal error request %w", userBA, err), http.StatusInternalServerError)
				return
			}
		}
		if !ok {
			api.HttpError(w, errors.New("login "+userBA+" not authorized"), http.StatusBadRequest)
			return
		}

		if next == nil {
			log.Fatal("handler requires next handler not nil")
		}
		//call further handler with context parameters
		ctx := context.WithValue(r.Context(), schema.CtxKeyUName, schema.CtxUName(userBA))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
