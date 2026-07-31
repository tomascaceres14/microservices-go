package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/shared/types"
)

const (
	BASE_OSRM_URL = "http://router.project-osrm.org"
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

func (s *TripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmAPIResponse, error) {
	url := fmt.Sprintf("%s/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", BASE_OSRM_URL,
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching OSRM API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading OSRM API response: %v", err)
	}

	var routeResponse *types.OsrmAPIResponse
	if err := json.Unmarshal(body, &routeResponse); err != nil {
		return nil, fmt.Errorf("error unparsing OSRM API response: %v", err)
	}

	routeResponse.PutoElQueLee = true

	return routeResponse, nil
}
