package grpc

import (
	"context"
	"ride-sharing/services/driver-service/internal/service"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
)

type GRPCHandler struct {
	pb.UnimplementedDriverServiceServer
	svc *service.DriverService
}

func NewGRPCHandler(server *grpc.Server, svc *service.DriverService) {
	pb.RegisterDriverServiceServer(server, &GRPCHandler{
		svc: svc,
	})

}

func (h *GRPCHandler) RegisterDriver(ctx context.Context, in *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	return nil, nil
}
func (h *GRPCHandler) UnregisterDriver(ctx context.Context, in *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	return nil, nil
}
