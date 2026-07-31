package service

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
)

type TripService struct {
	repo *repository.InmemRepository
}

func NewTripService(repo *repository.InmemRepository) *TripService {
	return &TripService{
		repo: repo,
	}
}

func (s *TripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {

	trip := domain.NewTripModel("ok", fare)
	return trip, nil
}
