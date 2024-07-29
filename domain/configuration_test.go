package domain_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
)

func TestNewDefaultPickUPWarehouse(t *testing.T) {
	retailerID, err := uuid.NewV4()
	assert.NoError(t, err)

	t.Run("success, the method works properly", func(t *testing.T) {
		got := domain.NewDefaultPickUPWarehouse(retailerID)
		assert.IsType(t, &domain.Configuration{}, got)
	})
}

func TestNewDefaultFailedPickUPWarehouse(t *testing.T) {
	retailerID, err := uuid.NewV4()
	assert.NoError(t, err)

	t.Run("success, the method works properly", func(t *testing.T) {
		got := domain.NewDefaultFailedPickUPWarehouse(retailerID)
		assert.IsType(t, &domain.Configuration{}, got)
	})
}

func TestNewDefaultFailedDeliveryWarehouse(t *testing.T) {
	retailerID, err := uuid.NewV4()
	assert.NoError(t, err)

	t.Run("success, the method works properly", func(t *testing.T) {
		got := domain.NewDefaultFailedDeliveryWarehouse(retailerID)
		assert.IsType(t, &domain.Configuration{}, got)
	})
}

func Test_RetailerConfiguration_AddDefaultConfigurations(t *testing.T) {
	t.Run("success, the method works properly", func(t *testing.T) {
		config := domain.RetailerConfiguration{}

		config.AddDefaultConfigurations()
		assert.Len(t, config.Configurations, 3)
	})
}
