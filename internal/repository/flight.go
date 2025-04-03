package repository

import (
	"context"

	"go-api/internal/apperr"
	"go-api/internal/model"
	"go-api/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type FlightRepository struct {
	db *pgx.Conn
}

func NewFlightRepository(db *pgx.Conn) usecase.FlightProvider {
	return &FlightRepository{db: db}
}

func (r *FlightRepository) GetAllFlights(ctx context.Context) ([]model.FlightDTO, error) {
	rows, err := r.db.Query(ctx, "SELECT id, destination_from, destination_to FROM public.flights WHERE deleted_at IS NULL")
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch flights")
	}
	defer rows.Close()

	var flights []model.FlightDTO
	for rows.Next() {
		var flight model.FlightDTO
		if err := rows.Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo); err != nil {
			return nil, errors.Wrap(err, "failed to scan flight row")
		}
		flights = append(flights, flight)
	}

	if len(flights) == 0 {
		return nil, apperr.ErrNotFound
	}

	return flights, nil
}

func (r *FlightRepository) GetFlightByID(ctx context.Context, id uuid.UUID) (*model.FlightDTO, error) {
	var flight model.FlightDTO
	err := r.db.QueryRow(ctx, "SELECT id, destination_from, destination_to FROM public.flights WHERE id=$1 AND deleted_at IS NULL", id).
		Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch flight by ID")
	}

	return &flight, nil
}

func (r *FlightRepository) InsertFlight(ctx context.Context, flight *model.FlightDTO) error {
	err := r.db.QueryRow(ctx,
		"INSERT INTO public.flights (destination_from, destination_to) VALUES ($1, $2) RETURNING id",
		flight.DestinationFrom, flight.DestinationTo).Scan(&flight.FlightID)

	if err != nil {
		return errors.Wrap(err, "failed to insert flight")
	}

	return nil
}

func (r *FlightRepository) UpdateFlight(ctx context.Context, flight *model.FlightDTO) error {
	result, err := r.db.Exec(ctx,
		"UPDATE public.flights SET destination_from=$1, destination_to=$2 WHERE id=$3 AND deleted_at IS NULL",
		flight.DestinationFrom, flight.DestinationTo, flight.FlightID)

	if err != nil {
		return errors.Wrap(err, "failed to update flight")
	}
	if result.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	return nil
}

func (r *FlightRepository) DeleteFlight(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, "UPDATE public.flights SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)

	if err != nil {
		return errors.Wrap(err, "failed to delete flight")
	}
	if result.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	return nil
}
