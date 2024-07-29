package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
)

func TestModel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       domain.Model
		wantErr bool
	}{
		{
			name:    "validate expected model",
			m:       domain.ModelStore,
			wantErr: false,
		},
		{
			name:    "validate not expected model",
			m:       domain.Model("invalid"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrUnsupportedRetailerModel, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestModel_String(t *testing.T) {
	t.Run("string representation for warehouse model", func(t *testing.T) {
		got := domain.ModelWarehouse.String()
		assert.Equal(t, "WAREHOUSE", got)
	})

	t.Run("string representation for store model", func(t *testing.T) {
		got := domain.ModelStore.String()
		assert.Equal(t, "STORE", got)
	})

	t.Run("string representation for not expected model", func(t *testing.T) {
		got := domain.Model("unexpected").String()
		assert.Equal(t, "unexpected", got)
	})
}
