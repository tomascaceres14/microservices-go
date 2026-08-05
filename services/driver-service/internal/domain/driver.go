package domain

import (
	"context"
	"ride-sharing/shared/types"
)

type Driver struct {
	Id, Name, ProfilePic, CarPlate, GeoHash, PackageSlug string
	Location                                             types.Coordinate
}

type DriverRepository interface {
	SaveDriver(ctx context.Context, driver *Driver) (*Driver, error)
	DeleteDriver(ctx context.Context, id string)
}

type DriverService interface {
	RegisterDriver(ctx context.Context, driver *Driver) *Driver
	UnregisterDriver(ctx context.Context, driver *Driver)
}
