package main

import (
	"log"
	"net/http"
	http_h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {
	log.Println("Starting Trip Service")

	inmemRepo := repository.NewInmemRepository()
	svc := service.NewTripService(inmemRepo)
	tripHandler := http_h.NewTripHandler(svc)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /preview", tripHandler.HandleCreateTrip)
	mux.HandleFunc("GET /preview", func(w http.ResponseWriter, r *http.Request) {
		res := "hola"
		w.Write([]byte(res))
	})

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP trip-service error: %v", err)
	}
}
