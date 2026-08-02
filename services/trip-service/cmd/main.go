package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"syscall"

	grpc_server "google.golang.org/grpc"
)

var (
	grpcAddr = env.GetString("HTTP_ADDR", ":9093")
)

func main() {

	// Initalize layers and server
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewTripService(inmemRepo)

	server := grpc_server.NewServer()
	grpc.NewGRPCHandler(server, svc)

	// Shutdown for K8S syscalls
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Run server
	log.Printf("Starting TRIP Service on port %s", grpcAddr)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("Error initializing TRIP Service: %s", err)
			cancel()
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down TRIP service")
	server.GracefulStop()
}
