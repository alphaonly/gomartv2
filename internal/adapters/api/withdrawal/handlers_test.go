package withdrawal_test

import (
	"bytes"
	"context"
	"fmt"

	wdHnd "github.com/alphaonly/gomartv2/internal/adapters/api/withdrawal"
	"github.com/alphaonly/gomartv2/internal/configuration"
	orderMocks "github.com/alphaonly/gomartv2/internal/mocks/order"
	mocks "github.com/alphaonly/gomartv2/internal/mocks/withdrawal"
	"github.com/alphaonly/gomartv2/internal/schema"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_PostWithdraw(t *testing.T) {

	type want struct {
		response    []byte
		status      int
		contentType string
	}

	type request struct {
		URL      string
		method   string
		testUser string
		order    string
		sum      float64
		body     []byte
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "test#1 positive",

			request: request{
				URL:      "/api/user/balance/withdraw",
				testUser: mocks.TestUser200,

				method: http.MethodPost,
				body:   mocks.TestJSON,
			},
			want: want{
				status: 200,
			},
		},
		{
			name: "test#2 positive - No funds",

			request: request{
				URL:      "/api/user/balance/withdraw",
				testUser: mocks.TestUser402,
				method:   http.MethodPost,
				body:     mocks.TestJSON,
			},
			want: want{
				status:      402,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "test#2 positive - Invalid order number",

			request: request{
				URL:      "/api/user/balance/withdraw",
				testUser: mocks.TestUser422,
				method:   http.MethodPost,
				body:     mocks.TestJSON,
			},
			want: want{
				status:      422,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	withdrawalStorage := mocks.NewWithdrawalStorage()
	withdrawalService := mocks.NewService()
	orderService := orderMocks.NewService()
	withdrawalHandler := wdHnd.NewHandler(withdrawalStorage, withdrawalService, orderService, cfg)
	// маршрутизация запросов обработчику
	rtr := NewRouter(withdrawalHandler)

	ts := httptest.NewServer(rtr)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(tt.request.method, tt.request.URL, bytes.NewBuffer(tt.request.body))
			req.WithContext(ctx)

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

func NewRouter(h wdHnd.Handler) chi.Router {

	var (
		withdraw = h.PostWithdraw()
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/api/user/balance/withdraw", withdraw)
	})

	return r
}
