package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"pinger/internal/config"
	"pinger/logger"
	"pinger/services"
	"time"
)

func main() {
	logger := logger.GetLogger()
	ctx := context.Background()
	router := httprouter.New()
	cfg := config.GetConfig()
	postgreSQLClient, err := services.NewClient(ctx, 3, cfg.Storage)
	if err != nil {
		os.Exit(1)
	}

	rep := services.NewRepository(logger, postgreSQLClient)
	hdlr := services.NewHandler(logger, postgreSQLClient, rep)
	hdlr.Register(router)
	start(logger, router)
}

func start(logger *logger.Logger, router *httprouter.Router) {
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	logger.Info("listen tcp")
	listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", "8088"))
	logger.Infof("server is listening port %s:%s", "127.0.0.1", "8088")

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
