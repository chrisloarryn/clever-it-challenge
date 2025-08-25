package inmemory

import (
	"context"
	"testing"

	"beers-challenge/internal/core/domain/beers"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	repo := NewRepository()
	beer := &beers.Beer{ID: 1, Name: "Test Beer"}

	err := repo.Save(context.Background(), beer)
	assert.NoError(t, err)

	savedBeer, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, beer, savedBeer)
}

func TestFindByIDNotFound(t *testing.T) {
	repo := NewRepository()
	_, err := repo.FindByID(context.Background(), 1)
	assert.Error(t, err)
}

func TestFindAll(t *testing.T) {
	repo := NewRepository()
	beer1 := &beers.Beer{ID: 1, Name: "Test Beer 1"}
	beer2 := &beers.Beer{ID: 2, Name: "Test Beer 2"}

	repo.Save(context.Background(), beer1)
	repo.Save(context.Background(), beer2)

	allBeers, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, allBeers, 2)
}

func TestExistsByID(t *testing.T) {
	repo := NewRepository()
	beer := &beers.Beer{ID: 1, Name: "Test Beer"}

	repo.Save(context.Background(), beer)

	exists, err := repo.ExistsByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByID(context.Background(), 2)
	assert.NoError(t, err)
	assert.False(t, exists)
}
