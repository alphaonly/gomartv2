package order_test

import (
	"bytes"
	"context"
	"github.com/alphaonly/gomartv2/internal/adapters/api/order"
	"github.com/alphaonly/gomartv2/internal/configuration"
	mocks "github.com/alphaonly/gomartv2/internal/mocks/order"
	userMocks "github.com/alphaonly/gomartv2/internal/mocks/user"
	"github.com/alphaonly/gomartv2/internal/schema"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// Example - an example to get order in mocked DB
func Example() {

	type request struct {
		URL      string
		method   string
		testUser string
		body     []byte
	}

	data := url.Values{}

	variants := []struct {
		name    string
		request request
	}{
		{
			name: "variant1",

			request: request{
				URL:      "/api/user/orders",
				method:   http.MethodGet,
				body:     []byte(""),
				testUser: "testuser",
			},
		},
	}
	// create a configuration
	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)
	// create an order mock storage
	orderStorage := mocks.NewOrderStorage()
	// create an order service
	orderService := mocks.NewService()
	//create a user mock service
	userService := userMocks.NewService()
	//create an order handlers
	orderHandler := order.NewHandler(orderStorage, orderService, userService, cfg)
	// routing request to handlers

	//a Simple route for getOrders
	var getOrders = orderHandler.GetOrders()

	//new chi router
	r := chi.NewRouter()

	//Apply route
	r.Route("/", func(r chi.Router) {
		r.Get("/api/user/orders", getOrders)
	})

	// a new http server
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, vv := range variants {

		req := httptest.NewRequest(vv.request.method, vv.request.URL, bytes.NewBufferString(data.Encode()))
		ctx := context.WithValue(req.Context(), schema.CtxKeyUName, schema.CtxUName(vv.request.testUser))
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		response := w.Result()
		if response.StatusCode != http.StatusOK {
			log.Fatalf("error code %v ", response.StatusCode)
		}

	}

}
