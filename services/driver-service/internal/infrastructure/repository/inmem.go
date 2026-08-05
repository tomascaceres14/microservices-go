package repository

import (
	"context"
	"ride-sharing/services/driver-service/internal/domain"
)

type InmemRepository struct {
	drivers map[string]*domain.Driver
}

func NewInmemRepository() *InmemRepository {
	return &InmemRepository{
		drivers: make(map[string]*domain.Driver),
	}
}

func (r *InmemRepository) SaveDriver(ctx context.Context, driver *domain.Driver) (*domain.Driver, error) {
	r.drivers[driver.Id] = driver
	return driver, nil
}

func (r *InmemRepository) DeleteDriver(ctx context.Context, id string) {
	delete(r.drivers, id)
}
