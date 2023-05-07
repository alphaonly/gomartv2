// Package api - contains help functions and common handlers for API part of entities
package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// BasicAuth - checks basic authentication whether it is valid
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// HTTPErrorW - forms http error message
func HTTPErrorW(w http.ResponseWriter, eStr string, err error, status int) {
	if err != nil {
		newE := fmt.Errorf(eStr+" %w", err)
		HTTPError(w, newE, status)
		log.Println("server:" + newE.Error())
	}
}

// HTTPError - another way of forming http error message
func HTTPError(w http.ResponseWriter, err error, status int) {
	if err != nil {
		http.Error(w, err.Error(), status)
		log.Println("server:" + err.Error())
	}
}

// GetPreviousParameter - Getting a context parameter from previous handler func in a chain
func GetPreviousParameter[T any, V any](r *http.Request, key V) (data T, err error) {
	var prev T
	var p any

	if p = r.Context().Value(key); p == nil {
		log.Println("got nil data from previous handler")
		return prev, errors.New("got nil data from previous handler")
	}

	return p.(T), nil

}
