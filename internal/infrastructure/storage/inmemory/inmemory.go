package inmemory

import (
	"context"
	"fmt"
	"sync"

	"beers-challenge/internal/core/domain/beers"
	"beers-challenge/internal/core/ports/secondary"
)

// Repository implements the secondary.BeerRepository interface for in-memory storage
type Repository struct {
	data map[int]*beers.Beer
	mu   sync.RWMutex
}

// NewRepository creates a new in-memory repository
func NewRepository() secondary.BeerRepository {
	return &Repository{
		data: make(map[int]*beers.Beer),
		mu:   sync.RWMutex{},
	}
}

// Save saves a beer to memory
func (r *Repository) Save(ctx context.Context, beer *beers.Beer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create a copy to avoid external modifications
	beerCopy := *beer
	r.data[beer.ID] = &beerCopy

	return nil
}

// FindByID finds a beer by its ID
func (r *Repository) FindByID(ctx context.Context, id int) (*beers.Beer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	beer, exists := r.data[id]
	if !exists {
		return nil, beers.NewDomainError("BEER_NOT_FOUND", fmt.Sprintf("Beer with ID %d not found", id), nil)
	}

	// Return a copy to avoid external modifications
	beerCopy := *beer
	return &beerCopy, nil
}

// FindAll finds all beers
func (r *Repository) FindAll(ctx context.Context) ([]beers.Beer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]beers.Beer, 0, len(r.data))
	for _, beer := range r.data {
		// Create copies to avoid external modifications
		beerCopy := *beer
		result = append(result, beerCopy)
	}

	return result, nil
}

// ExistsByID checks if a beer exists by its ID
func (r *Repository) ExistsByID(ctx context.Context, id int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.data[id]
	return exists, nil
}
