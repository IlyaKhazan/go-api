package usecase

import (
	"context"

	"go-api/internal/model"
	"go-api/internal/repository"

	uuid "github.com/satori/go.uuid"
)

type FlightUsecase struct {
	flightUC FlightProvider
}

func NewFlightUsecase(repo *repository.FlightRepository) FlightProvider {
	return &FlightUsecase{flightUC: repo}
}

func (uc *FlightUsecase) GetAllFlights(ctx context.Context) ([]model.FlightDTO, error) {
	return uc.flightUC.GetAllFlights(ctx)
}

func (uc *FlightUsecase) GetFlightByID(ctx context.Context, id uuid.UUID) (*model.FlightDTO, error) {
	return uc.flightUC.GetFlightByID(ctx, id)
}

func (uc *FlightUsecase) InsertFlight(ctx context.Context, flight *model.FlightDTO) error {
	return uc.flightUC.InsertFlight(ctx, flight)
}

func (uc *FlightUsecase) UpdateFlight(ctx context.Context, flight *model.FlightDTO) error {
	return uc.flightUC.UpdateFlight(ctx, flight)
}

func (uc *FlightUsecase) DeleteFlight(ctx context.Context, id uuid.UUID) error {
	return uc.flightUC.DeleteFlight(ctx, id)
}
