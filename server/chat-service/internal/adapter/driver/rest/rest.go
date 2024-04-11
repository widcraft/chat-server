package rest

import (
	"chat-service/internal/adapter/driven/config"
	"chat-service/internal/port/driven"
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Register interface {
	Register(router gin.IRoutes)
}

type Rest struct {
	logger driven.Logger
	server *http.Server
}

func New(configStore *config.Store, logger driven.Logger) *Rest {
	return &Rest{
		logger: logger,
		server: &http.Server{
			Handler:      gin.Default(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			Addr:         ":" + configStore.GetRestPort(),
		},
	}
}

func (r *Rest) Register(handlers ...Register) {
	group := r.server.Handler.(*gin.Engine).Group("/api/v1")

	for _, handler := range handlers {
		handler.Register(group)
	}
}

func (r *Rest) Run(_ context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	if err := r.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (r *Rest) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return r.server.Shutdown(ctx)
}
