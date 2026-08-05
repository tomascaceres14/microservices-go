package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/driver-service/internal/infrastructure/grpc"
	"ride-sharing/services/driver-service/internal/infrastructure/repository"
	"ride-sharing/services/driver-service/internal/service"
	"ride-sharing/shared/env"
	"syscall"

	grpc_server "google.golang.org/grpc"
)

var (
	grpcAddr = env.GetString("GRPC_ADDR", ":9092")
)

func main() {

	server := grpc_server.NewServer()
	repo := repository.NewInmemRepository()
	svc := service.NewDriverService(repo)
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
	log.Printf("Starting DRIVER Service on port %s", grpcAddr)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("Error initializing DRIVER Service: %s", err)
			cancel()
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down DRIVER service")
	server.GracefulStop()
}
