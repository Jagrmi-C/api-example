package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
)

func TestCountry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       domain.Country
		wantErr bool
	}{
		{
			name:    "success, validate supported country: ES",
			m:       domain.CountryES,
			wantErr: false,
		},
		{
			name:    "success, validate supported country: IT",
			m:       domain.CountryIT,
			wantErr: false,
		},
		{
			name:    "failed, validate not supported country: PL",
			m:       domain.Country("PL"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrUnsupportedCountry, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestCountry_String(t *testing.T) {
	t.Run("success, the method works properly", func(t *testing.T) {
		got := domain.CountryES.String()
		assert.Equal(t, "ES", got)
	})
}
