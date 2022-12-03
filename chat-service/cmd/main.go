package main

import (
	"net"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/widcraft/chat-service/external/workerpool"
	"github.com/widcraft/chat-service/internal/adapter/repository"
	"github.com/widcraft/chat-service/internal/adapter/repository/redis"
	"github.com/widcraft/chat-service/internal/adapter/websocket"
	"github.com/widcraft/chat-service/internal/application"
)

var logger = log.New()
var redisDb *redis.Redis
var (
	wg       = &sync.WaitGroup{}
	chatPool = workerpool.New(wg, 1)
)

var (
	ws *websocket.Websocket
)

// env
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}
}

// DB
func init() {
	var err error
	redisDb, err = redis.New(logger, net.JoinHostPort(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), os.Getenv("REDIS_PASSWORD"), 0)

	if err != nil {
		shutdown(redisDb)
		logger.Fatalf("redis connection failure: %s", err)
	}
}

// servers
func init() {
	chatRepo := repository.NewChatRepository(logger, redisDb)
	chatApp := application.NewChatApp(logger, chatPool, chatRepo)
	ws = websocket.New(logger, chatApp)
}

func main() {
	defer gracefulShutdown()
	go ws.Run(os.Getenv("WS_PORT"))
}

func gracefulShutdown() {
	defer shutdown(chatPool, ws, redisDb)

	terminationChan := make(chan os.Signal, 1)
	signal.Notify(terminationChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-terminationChan
	wg.Wait()
}

func shutdown(targets ...interface{ Close() }) {
	for _, target := range targets {
		if !reflect.ValueOf(target).IsNil() {
			target.Close()
		}
	}
}
