package service

import (
	"context"
	"ride-sharing/services/driver-service/internal/domain"
)

type DriverService struct {
	repo domain.DriverRepository
}

func NewDriverService(r domain.DriverRepository) *DriverService {
	return &DriverService{
		repo: r,
	}
}

func (s *DriverService) RegisterDriver(ctx context.Context, driver *domain.Driver) *domain.Driver {
	return nil
}
func (s *DriverService) UnregisterDriver(ctx context.Context, driver *domain.Driver) {}
