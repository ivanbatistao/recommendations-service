package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivanbatistao/recommendation-service/configs"
	httpgin "github.com/ivanbatistao/recommendation-service/internal/infrastructure/http/gin"
)

func main() {
	config := configs.Load()

	server := httpgin.NewServer(config.Port)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.Start()
	}()

	log.Printf("server running on :%s", config.Port)

	shutdownSignal := make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignal,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		log.Fatalf("server error: %v", err)

	case signal := <-shutdownSignal:
		log.Printf("shutdown signal received: %s", signal)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
