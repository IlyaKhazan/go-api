package usecase

import (
	"context"

	"go-api/internal/model"
)

type FlightProvider interface {
	GetAllFlights(ctx context.Context) ([]model.FlightDTO, error)
	GetFlightByID(ctx context.Context, id int) (*model.FlightDTO, error)
	InsertFlight(ctx context.Context, flight *model.FlightDTO) error
	UpdateFlight(ctx context.Context, flight *model.FlightDTO) error
	DeleteFlight(ctx context.Context, id int) error
}
