// order_test - пакет тестирования хэндоеров для обработки заказов
package order_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alphaonly/gomartv2/internal/adapters/api/order"
	"github.com/alphaonly/gomartv2/internal/schema"

	"github.com/alphaonly/gomartv2/internal/configuration"
	mocks "github.com/alphaonly/gomartv2/internal/mocks/order"
	userMocks "github.com/alphaonly/gomartv2/internal/mocks/user"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandler_GetOrders(t *testing.T) {

	type want struct {
		response    []byte
		status      int
		contentType string
	}

	type request struct {
		URL      string
		method   string
		testUser string
		body     []byte
	}

	data := url.Values{}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "test#1 positive",

			request: request{
				URL:      "/api/user/orders",
				method:   http.MethodGet,
				body:     []byte(""),
				testUser: "testuser",
			},
			want: want{
				status:      200,
				contentType: "application/json",
			},
		},
		{
			name: "test#2 negative - no user's orders",

			request: request{
				URL:      "/api/user/orders",
				method:   http.MethodGet,
				body:     []byte(""),
				testUser: "testuser2",
			},
			want: want{
				status:      204,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	orderStorage := mocks.NewOrderStorage()
	orderService := mocks.NewService()
	userService := userMocks.NewService()

	orderHandler := order.NewHandler(orderStorage, orderService, userService, cfg)
	// маршрутизация запросов обработчику
	rtr := NewRouter(orderHandler)

	ts := httptest.NewServer(rtr)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(tt.request.method, tt.request.URL, bytes.NewBufferString(data.Encode()))
			ctx := context.WithValue(req.Context(), schema.CtxKeyUName, schema.CtxUName(tt.request.testUser))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			rtr.ServeHTTP(w, req)

			response := w.Result()
			if response.StatusCode != tt.want.status {
				t.Errorf("error code %v want %v", response.StatusCode, tt.want.status)
				fmt.Println(response)
				fmt.Println(w.Body.String())

			}

			if response.Header.Get("Content-type") != tt.want.contentType {
				t.Errorf("error contentType %v want %v", response.Header.Get("Content-type"), tt.want.contentType)
			}
			err := response.Body.Close()
			if err != nil {
				t.Errorf("response body close error: %v response", response.Body)
			}
		})

	}

}

func NewRouter(h order.Handler) chi.Router {

	var (
		getOrders = h.GetOrders()
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/api/user/orders", getOrders)
	})

	return r
}
