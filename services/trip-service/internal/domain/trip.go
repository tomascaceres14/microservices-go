package domain

import (
	"context"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID `json:"id"`
	UserID   string             `json:"userID"`
	Status   string             `json:"status"`
	RideFare *RideFareModel     `json:"ride_fare"`
}

func NewTripModel(status string, fare *RideFareModel) *TripModel {
	return &TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   status,
		RideFare: fare,
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmAPIResponse, error)
}
