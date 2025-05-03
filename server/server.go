package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michaelyusak/go-helper/helper"
	"github.com/michaelyusak/xyz-kredit-plus/config"
)

func Init() {
	log := helper.NewLogrus()

	config := config.Init(log)

	router := createRouter(config, log)

	srv := http.Server{
		Handler: router,
		Addr:    config.Port,
	}

	go func() {
		log.Infof("Sever running on port: %s", config.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 10)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Server shutdown gracefully ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GracefulPeriod)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown: %s", err.Error())
	}

	<-ctx.Done()

	log.Infof("Timeout of %v seconds", config.GracefulPeriod)
	log.Info("Server exited")
}
