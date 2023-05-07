// Package user - this API part of user entity that contains user' handlers functionality.
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
	Register() http.HandlerFunc                   // a function to implement HTTP POST user data handler to save it
	Login() http.HandlerFunc                      // a function to implement HTTP POST user data handler to login user
	BasicAuth(next http.Handler) http.HandlerFunc // a technical handler to check user's basic authentication
}

type handler struct {
	Storage       user.Storage                       // a pointer to user storage interface implementation
	Service       user.Service                       // a pointer to user service interface implementation
	Configuration *configuration.ServerConfiguration // a pointer to a server configuration
}

// NewHandler - it is a factory that returns an instance of user's Handler implementation.
func NewHandler(storage user.Storage, service user.Service, configuration *configuration.ServerConfiguration) Handler {
	return &handler{
		Storage:       storage,
		Service:       service,
		Configuration: configuration,
	}
}

// Register - a handler func to catch POST request of adding new user(registration)
func (h *handler) Register() http.HandlerFunc {
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
			if errors.Is(err, user.ErrUserPassEmpty) {
				http.Error(w, "login "+u.User+": bad request", http.StatusBadRequest)
				return
			}
			if errors.Is(err, user.ErrLoginOccupied) {
				http.Error(w, "login "+u.User+"is occupied", http.StatusConflict)
				return
			}
			if errors.Is(err, user.ErrInternal) {
				http.Error(w, "login "+u.User+"register internal error", http.StatusInternalServerError)
				return
			}
		}
		//Response
		log.Printf("Respond in header basic authorization: user:%v password: %v", u.User, u.Password)
		w.Header().Add("Authorization", "Basic "+api.BasicAuth(u.User, u.Password))
		w.WriteHeader(http.StatusOK)

	}
}

// Login - a handler func to catch POST request of logging in a new user
func (h *handler) Login() http.HandlerFunc {
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
			if errors.Is(err, user.ErrUserPassEmpty) {
				http.Error(w, "login "+u.User+": bad request", http.StatusBadRequest)
				return
			}
			if errors.Is(err, user.ErrLogPassUnknown) {
				api.HTTPErrorW(w, "authorization error", err, http.StatusUnauthorized)
				return
			}
			if errors.Is(err, user.ErrLoginOccupied) {
				http.Error(w, "login "+u.User+"is occupied", http.StatusConflict)
				return
			}
			api.HTTPErrorW(w, "login "+u.User+"register internal error", err, http.StatusInternalServerError)
			return
		}
		//Response
		log.Printf("Respond in header basic authorization: user:%v password: %v", u.User, u.Password)
		w.Header().Add("Authorization", "Basic "+api.BasicAuth(u.User, u.Password))
		w.WriteHeader(http.StatusOK)
	}
}

// BasicAuth - a handler func to check basic authentication of user
func (h *handler) BasicAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("BasicUserAuthorization invoked")
		//Basic authentication
		userBA, passBA, ok := r.BasicAuth()
		if !ok {
			api.HTTPError(w, fmt.Errorf("basic authentication is not ok"), http.StatusUnauthorized)
			return
		}
		log.Printf("basic authorization check: user: %v, password: %v", userBA, passBA)

		var err error
		ok, err = h.Service.CheckIfUserAuthorized(r.Context(), userBA, passBA)
		if err != nil {
			if strings.Contains(err.Error(), "400") {
				api.HTTPError(w, fmt.Errorf("login %v: bad request %w", userBA, err), http.StatusBadRequest)
				return
			}
			if strings.Contains(err.Error(), "500") {
				api.HTTPError(w, fmt.Errorf("login %v: server internal error request %w", userBA, err), http.StatusInternalServerError)
				return
			}
		}
		if !ok {
			api.HTTPError(w, errors.New("login "+userBA+" not authorized"), http.StatusBadRequest)
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
