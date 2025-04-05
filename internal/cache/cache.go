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
	flightsRepo repository.FlightProvider
	flights     map[uuid.UUID]*model.FlightDTO
	mu          sync.RWMutex
}

func NewCacheDecorator(flightsRepo repository.FlightProvider) repository.FlightProvider {
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
	c.mu.RLock()
	if len(c.flights) > 0 {
		flights := make([]model.FlightDTO, 0, len(c.flights))
		for _, flight := range c.flights {
			flights = append(flights, *flight)
		}
		c.mu.RUnlock()
		return flights, nil
	}
	c.mu.RUnlock()

	flightsFromRepo, err := c.flightsRepo.GetAllFlights(ctx)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	for _, f := range flightsFromRepo {
		flight := f
		c.flights[f.FlightID] = &flight
	}
	c.mu.Unlock()

	return flightsFromRepo, nil
}

func (c *Decorator) GetFlightByID(ctx context.Context, id uuid.UUID) (*model.FlightDTO, error) {
	if flight, exists := c.Get(id); exists {
		return flight, nil
	}

	if flight, err := c.flightsRepo.GetFlightByID(ctx, id); err == nil {
		c.Set(flight)
		return flight, nil
	} else {
		return nil, err
	}
}

func (c *Decorator) InsertFlight(_ context.Context, flight *model.FlightDTO) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.flights[flight.FlightID]; exists {
		return errors.New("flight already exists")
	}
	c.flights[flight.FlightID] = flight
	return nil
}

func (c *Decorator) UpdateFlight(_ context.Context, flight *model.FlightDTO) error {
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
