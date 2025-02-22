package usecase

import (
	"context"
)

type FlightProvider interface {
	GetAllFlights(ctx context.Context) ([]Flight, error)
	GetFlightByID(ctx context.Context, id int) (*Flight, error)
	InsertFlight(ctx context.Context, flight *Flight) error
	UpdateFlight(ctx context.Context, flight *Flight) error
	DeleteFlight(ctx context.Context, id int) error
}
