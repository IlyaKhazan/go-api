package cache

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"go-api/internal/model"
	"go-api/internal/repository"

	uuid "github.com/satori/go.uuid"
)

type Decorator struct {
	mu          sync.RWMutex
	flights     map[uuid.UUID]flightWithTTL
	flightsRepo repository.FlightProvider
	ttl         time.Duration
}

type flightWithTTL struct {
	data      *model.FlightDTO
	expiresAt time.Time
}

func NewDecorator(flightsRepo repository.FlightProvider, ttl time.Duration) *Decorator {
	return &Decorator{flightsRepo: flightsRepo,
		flights: make(map[uuid.UUID]flightWithTTL),
		ttl:     ttl}
}

func (c *Decorator) Get(id uuid.UUID) (*model.FlightDTO, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	flight, exists := c.flights[id]
	if !exists || time.Now().After(flight.expiresAt) {
		return nil, false
	}

	return flight.data, true
}

func (c *Decorator) Set(flight *model.FlightDTO) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flights[flight.FlightID] = flightWithTTL{
		data:      flight,
		expiresAt: time.Now().Add(c.ttl),
	}
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
	c.flights[flight.FlightID] = flightWithTTL{
		data:      flight,
		expiresAt: time.Now().Add(c.ttl),
	}
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

	c.flights[flight.FlightID] = flightWithTTL{
		data:      flight,
		expiresAt: time.Now().Add(c.ttl),
	}
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

func (c *Decorator) StartCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.mu.Lock()
			now := time.Now()
			removed := 0

			for id, flight := range c.flights {
				if now.After(flight.expiresAt) {
					delete(c.flights, id)
					removed++
				}
			}

			c.mu.Unlock()
			slog.Info("✅ cache cleanup cycle complete", "expired_removed", removed, "remaining", len(c.flights))
		}
	}()
}
