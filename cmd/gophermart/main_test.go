package main_test

import (
	"context"
	"github.com/alphaonly/gomartv2/internal/adapters/api"
	"github.com/alphaonly/gomartv2/internal/adapters/api/router"
	"github.com/alphaonly/gomartv2/internal/composites"
	"github.com/alphaonly/gomartv2/internal/configuration"
	"github.com/alphaonly/gomartv2/internal/pkg/common/logging"
	"github.com/alphaonly/gomartv2/internal/pkg/dbclient/postgres"
	"github.com/alphaonly/gomartv2/internal/pkg/server"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {

	tests := []struct {
		name string
		URL  string
		want string
	}{
		{
			name: "test#1 - Positive: server accessible",
			URL:  "http://localhost:8080/check",
			want: "200 OK",
		},
		{
			name: "test#2 - Negative: server do not respond",
			URL:  "http://localhost:8080/chek",
			want: "404 Not Found",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	dbClient := postgres.NewPostgresClient(ctx, cfg.DatabaseURI)

	UserComposite := composites.NewUserComposite(dbClient, cfg)
	OrderComposite := composites.NewOrderComposite(dbClient, cfg)
	WithdrawalComposite := composites.NewWithdrawalComposite(dbClient, cfg, UserComposite.Storage, OrderComposite.Service)

	handlerComposite := composites.NewHandlerComposite(
		api.NewHandler(cfg),
		UserComposite.Handler,
		OrderComposite.Handler,
		WithdrawalComposite.Handler,
	)

	// маршрутизация запросов обработчику
	rtr := router.NewRouter(handlerComposite)

	httpServer := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: rtr,
	}

	srv := server.NewServer(httpServer)

	go srv.Run()

	//resty client
	keys := make(map[string]string)
	keys["Content-Type"] = "plain/text"
	keys["Accept"] = "plain/text"

	client := resty.New()

	r := client.R().
		SetHeaders(keys)

	for _, tt := range tests {

		t.Run(tt.name, func(tst *testing.T) {
			//Up server for 3 seconds

			//wait for server is up
			time.Sleep(time.Second * 2)

			resp, err := r.Get(tt.URL)
			if err != nil {
				t.Logf("send new request error:%v", err)
			}
			t.Logf("get returned status:%v", resp.Status())
			if !assert.Equal(t, tt.want, resp.Status()) {
				t.Error("Server responded unexpectedly")

			}

		})
	}

	err := srv.Stop(ctx)
	logging.LogFatal(err)

}
