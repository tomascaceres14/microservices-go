package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	http_h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"syscall"
	"time"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {
	log.Println("Starting TRIP Service")

	inmemRepo := repository.NewInmemRepository()
	svc := service.NewTripService(inmemRepo)
	tripHandler := http_h.NewTripHandler(svc)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /preview", tripHandler.HandleCreateTrip)
	mux.HandleFunc("GET /preview", func(w http.ResponseWriter, r *http.Request) {
		res := "hola"
		w.Write([]byte(res))
	})

	sv := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	svStartCh := make(chan error, 1)
	go func() {
		log.Printf("TRIP Service listening on port %s", httpAddr)
		svStartCh <- sv.ListenAndServe()
	}()

	// Graceful shutdown. Wait for OS, K8s signal or error and shutdown after max 10s
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-svStartCh:
		log.Printf("Error starting the server: %s", err)

	case sig := <-shutdown:
		log.Printf("Server shutting down with signal: %s", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		if err := sv.Shutdown(ctx); err != nil {
			log.Printf("Could not shutdown gracefully: %s", err)
			sv.Close()
		}
	}
}
