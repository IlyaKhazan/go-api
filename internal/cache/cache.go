package cache

import (
	"context"
	"errors"
	"sync"

	"go-api/internal/model"
	"go-api/internal/repository"

	uuid "github.com/satori/go.uuid"
)

type Decorator struct {
	mu          sync.RWMutex
	flights     map[uuid.UUID]*model.FlightDTO
	flightsRepo repository.FlightProvider
}

func NewDecorator(flightsRepo repository.FlightProvider) repository.FlightProvider {
	return &Decorator{flightsRepo: flightsRepo}
}

func (c *Decorator) Get(id uuid.UUID) (*model.FlightDTO, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	flight, exists := c.flights[id]
	return flight, exists
}

func (c *Decorator) Set(flight *model.FlightDTO) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flights[flight.FlightID] = flight
}

func (c *Decorator) GetAllFlights(ctx context.Context) ([]model.FlightDTO, error) {
	return c.flightsRepo.GetAllFlights(ctx)
}

func (c *Decorator) GetFlightByID(ctx context.Context, id uuid.UUID) (*model.FlightDTO, error) {
	if flight, exists := c.Get(id); exists {
		return flight, nil
	}

	flight, err := c.flightsRepo.GetFlightByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Set(flight)
	return flight, nil
}

func (c *Decorator) InsertFlight(_ context.Context, flight *model.FlightDTO) error {
	if err := c.flightsRepo.InsertFlight(context.Background(), flight); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.flights[flight.FlightID]; exists {
		return errors.New("flight already exists")
	}
	c.flights[flight.FlightID] = flight
	return nil
}

func (c *Decorator) UpdateFlight(_ context.Context, flight *model.FlightDTO) error {
	if err := c.flightsRepo.UpdateFlight(context.Background(), flight); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.flights[flight.FlightID]; !exists {
		return errors.New("flight not found in cache")
	}

	c.flights[flight.FlightID] = flight
	return nil
}

func (c *Decorator) DeleteFlight(ctx context.Context, id uuid.UUID) error {
	if err := c.flightsRepo.DeleteFlight(ctx, id); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.flights, id)
	return nil
}
