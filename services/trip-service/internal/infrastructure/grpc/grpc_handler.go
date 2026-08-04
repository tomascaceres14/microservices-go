package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	pb.UnimplementedTripServiceServer
	svc domain.TripService
}

func NewGRPCHandler(server *grpc.Server, svc domain.TripService) *GRPCHandler {
	handler := &GRPCHandler{
		svc: svc,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *GRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {

	pickup := &types.Coordinate{
		Latitude:  req.GetStartLocation().GetLatitude(),
		Longitude: req.GetStartLocation().GetLongitude(),
	}

	destination := &types.Coordinate{
		Latitude:  req.GetEndLocation().GetLatitude(),
		Longitude: req.GetEndLocation().GetLongitude(),
	}

	route, err := h.svc.GetRoute(ctx, pickup, destination)
	if err != nil {
		log.Printf("Error fetching route: %s", err)
		return nil, status.Errorf(codes.Internal, "Error fetching route: %s", err)
	}

	estimatedFares := h.svc.EstimatePackagesPriceWithRoute(route)
	fares, err := h.svc.GenerateTripFares(ctx, estimatedFares, req.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error saving fares: %s", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}, nil

}

func (h *GRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fare, err := h.svc.GetAndValidateFare(ctx, req.GetRideFareID(), req.GetUserID())
	if err != nil {
		return &pb.CreateTripResponse{}, status.Errorf(codes.Internal, "Error fetching fare: %s", err)
	}

	trip, err := h.svc.CreateTrip(ctx, fare)
	if err != nil {
		return &pb.CreateTripResponse{}, status.Errorf(codes.Internal, "Error fetching trip: %s", err)
	}

	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil

}

func (h *GRPCHandler) mustEmbedUnimplementedTripServiceServer() {}
