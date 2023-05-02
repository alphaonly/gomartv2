package user_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alphaonly/gomartv2/internal/adapters/api/user"
	UserHnd "github.com/alphaonly/gomartv2/internal/adapters/api/user"
	"github.com/alphaonly/gomartv2/internal/configuration"
	mocks "github.com/alphaonly/gomartv2/internal/mocks/user"
	"github.com/alphaonly/gomartv2/internal/schema"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Register(t *testing.T) {

	type want struct {
		response    []byte
		status      int
		contentType string
	}

	type request struct {
		URL      string
		method   string
		testUser string
		testPass string
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
				URL:      "/api/user/register",
				testUser: mocks.TestUser200,
				testPass: mocks.TestPass200,
				method:   http.MethodPost,
				body:     []byte(fmt.Sprintf(`{"login":"%v","password": "%v"}`, mocks.TestUser200, mocks.TestPass200)),
			},
			want: want{
				status: 200,
			},
		},
		{
			name: "test#2 negative - Login Occupied 409 Error",
			request: request{
				URL:      "/api/user/register",
				method:   http.MethodPost,
				body:     []byte(fmt.Sprintf(`{"login":"%v","password": "%v"}`, mocks.TestUser409, mocks.TestPass409)),
				testUser: mocks.TestUser409,
				testPass: mocks.TestPass409,
			},
			want: want{
				status:      409,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "test#2 negative - Intenal Error",
			request: request{
				URL:      "/api/user/register",
				method:   http.MethodPost,
				body:     []byte(fmt.Sprintf(`{"login":"%v","password": "%v"}`, mocks.TestUser500, mocks.TestPass500)),
				testUser: mocks.TestUser500,
				testPass: mocks.TestPass500,
			},
			want: want{
				status:      500,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	userStorage := mocks.NewUserStorage()
	userService := mocks.NewService()

	userHandler := user.NewHandler(userStorage, userService, cfg)
	// маршрутизация запросов обработчику
	rtr := NewRouter(userHandler)

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

func NewRouter(h UserHnd.Handler) chi.Router {

	var (
		register = h.Register()
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/api/user/register", register)
	})

	return r
}
