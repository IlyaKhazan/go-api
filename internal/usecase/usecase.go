package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type FlightUsecase struct {
	conn *pgx.Conn
}

func New(db *pgx.Conn) FlightProvider {
	return &FlightUsecase{conn: db}
}

type Flight struct {
	FlightID        int    `json:"id"`
	DestinationFrom string `json:"destination_from"`
	DestinationTo   string `json:"destination_to"`
}

func (uc *FlightUsecase) GetAllFlights(ctx context.Context) ([]Flight, error) {
	rows, err := uc.conn.Query(ctx, "SELECT id, destination_from, destination_to FROM public.flights WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flights []Flight
	for rows.Next() {
		var flight Flight
		if err := rows.Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo); err != nil {
			return nil, err
		}
		flights = append(flights, flight)
	}

	if len(flights) == 0 {
		return nil, pgx.ErrNoRows
	}

	return flights, nil
}

func (uc *FlightUsecase) GetFlightByID(ctx context.Context, id int) (*Flight, error) {
	var flight Flight
	err := uc.conn.QueryRow(ctx, "SELECT id, destination_from, destination_to FROM public.flights WHERE id=$1 AND deleted_at IS NULL", id).
		Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}

	return &flight, err
}

func (uc *FlightUsecase) InsertFlight(ctx context.Context, flight *Flight) error {
	return uc.conn.QueryRow(ctx,
		"INSERT INTO public.flights (destination_from, destination_to) VALUES ($1, $2) RETURNING id",
		flight.DestinationFrom, flight.DestinationTo).Scan(&flight.FlightID)
}

func (uc *FlightUsecase) UpdateFlight(ctx context.Context, flight *Flight) error {
	result, err := uc.conn.Exec(ctx,
		"UPDATE public.flights SET destination_from=$1, destination_to=$2 WHERE id=$3 AND deleted_at IS NULL",
		flight.DestinationFrom, flight.DestinationTo, flight.FlightID)

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}

func (uc *FlightUsecase) DeleteFlight(ctx context.Context, id int) error {
	result, err := uc.conn.Exec(ctx, "UPDATE public.flights SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return err
}
