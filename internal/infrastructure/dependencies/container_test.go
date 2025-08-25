package dependencies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer(t *testing.T) {
	container, err := NewContainer()
	assert.NoError(t, err)
	assert.NotNil(t, container)

	assert.NotNil(t, container.GetConfig())
	assert.NotNil(t, container.GetLogger())
	assert.NotNil(t, container.GetBeerRepository())
	assert.NotNil(t, container.GetCurrencyService())
	assert.NotNil(t, container.GetBeerService())
	assert.NotNil(t, container.GetHTTPServer())

	err = container.Close()
	assert.NoError(t, err)
}

func TestGetters(t *testing.T) {
	container, _ := NewContainer()

	assert.Equal(t, container.config, container.GetConfig())
	assert.Equal(t, container.logger, container.GetLogger())
	assert.Equal(t, container.beerRepository, container.GetBeerRepository())
	assert.Equal(t, container.currencyService, container.GetCurrencyService())
	assert.Equal(t, container.beerService, container.GetBeerService())
	assert.Equal(t, container.httpServer, container.GetHTTPServer())
}
