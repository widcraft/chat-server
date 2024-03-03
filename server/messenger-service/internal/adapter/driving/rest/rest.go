package rest

import (
	"context"
	"errors"
	"messenger-service/internal/adapter/driving/rest/chat"
	"messenger-service/internal/port/driven"
	"messenger-service/internal/port/driving"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	logger           driven.LoggerPort
	server           *http.Server
	messengerService driving.MessengerServicePort
}

func New(logger driven.LoggerPort, messengerService driving.MessengerServicePort) *Rest {
	router := gin.Default()
	group := router.Group("/api/v1")

	chat.NewHandler(logger, messengerService).Register(group)

	return &Rest{
		logger:           logger,
		messengerService: messengerService,
		server: &http.Server{
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (rest *Rest) Run(port string) {
	rest.server.Addr = ":" + port
	err := rest.server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		rest.logger.Errorf("websocket server error: %s", err)
	}
}

func (rest *Rest) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return rest.server.Shutdown(ctx)
}
