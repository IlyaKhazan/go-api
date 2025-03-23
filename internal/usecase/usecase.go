package usecase

import (
	"context"

	"go-api/internal/model"
	"go-api/internal/repository"

	uuid "github.com/satori/go.uuid"
)

type FlightUsecase struct {
	repo FlightProvider
}

func NewFlightUsecase(repo *repository.FlightRepository) FlightProvider {
	return &FlightUsecase{repo: repo}
}

func (uc *FlightUsecase) GetAllFlights(ctx context.Context) ([]model.FlightDTO, error) {
	return uc.repo.GetAllFlights(ctx)
}

func (uc *FlightUsecase) GetFlightByID(ctx context.Context, id uuid.UUID) (*model.FlightDTO, error) {
	return uc.repo.GetFlightByID(ctx, id)
}

func (uc *FlightUsecase) InsertFlight(ctx context.Context, flight *model.FlightDTO) error {
	return uc.repo.InsertFlight(ctx, flight)
}

func (uc *FlightUsecase) UpdateFlight(ctx context.Context, flight *model.FlightDTO) error {
	return uc.repo.UpdateFlight(ctx, flight)
}

func (uc *FlightUsecase) DeleteFlight(ctx context.Context, id uuid.UUID) error {
	return uc.repo.DeleteFlight(ctx, id)
}
