package main

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewTripService(inmemRepo)
	fare := &domain.RideFareModel{UserID: "4343"}
	t, err := svc.CreateTrip(context.Background(), fare)
	if err != nil {
		print(err)
	}

	log.Println(t)

	for {
		time.Sleep(1 * time.Second)
	}
}
