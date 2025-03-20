package usecase

import (
	"context"

	"go-api/internal/model"

	uuid "github.com/satori/go.uuid"
)

type FlightProvider interface {
	GetAllFlights(ctx context.Context) ([]model.FlightDTO, error)
	GetFlightByID(ctx context.Context, id uuid.UUID) (*model.FlightDTO, error)
	InsertFlight(ctx context.Context, flight *model.FlightDTO) error
	UpdateFlight(ctx context.Context, flight *model.FlightDTO) error
	DeleteFlight(ctx context.Context, id uuid.UUID) error
}
