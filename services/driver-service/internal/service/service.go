package service

import (
	"context"
	"math/rand/v2"
	"ride-sharing/services/driver-service/internal/domain"
	driver_utils "ride-sharing/services/driver-service/utils"
	"ride-sharing/shared/types"
	"ride-sharing/shared/util"
	"sync"
)

type DriverService struct {
	repo domain.DriverRepository
	mu   sync.RWMutex
}

func NewDriverService(r domain.DriverRepository) *DriverService {
	return &DriverService{
		repo: r,
	}
}

func (s *DriverService) RegisterDriver(ctx context.Context, id, packageSlug string) (*domain.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := rand.IntN(10)
	route := driver_utils.PredefinedRoutes[randomIndex]
	avatar := util.GetRandomAvatar(randomIndex)
	plate := driver_utils.GenerateRandomPlate()

	driver := &domain.Driver{
		Id:             id,
		Name:           "Susana Gimenez",
		ProfilePicture: avatar,
		CarPlate:       plate,
		GeoHash:        "geohash",
		PackageSlug:    packageSlug,
		Location: types.Coordinate{
			Latitude:  route[0][0],
			Longitude: route[0][1],
		},
	}
	return s.repo.SaveDriver(ctx, driver)
}

func (s *DriverService) UnregisterDriver(ctx context.Context, id string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.repo.DeleteDriver(ctx, id)
}
