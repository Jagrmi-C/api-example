package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
)

func TestDeliveryType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       domain.DeliveryType
		wantErr bool
	}{
		{
			name:    "DIRECT",
			m:       domain.DeliveryTypeDirect,
			wantErr: false,
		},
		{
			name:    "REVERSE",
			m:       domain.DeliveryTypeReverse,
			wantErr: false,
		},
		{
			name:    "FAKE",
			m:       domain.DeliveryType("FAKE"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestDeliveryType_String(t *testing.T) {
	t.Run("happy path, the method works properly", func(t *testing.T) {
		got := domain.DeliveryTypeDirect.String()
		assert.Equal(t, "DIRECT", got)
	})
}
