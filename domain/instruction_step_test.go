package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
)

func TestStepName_Validate(t *testing.T) {
	tests := []struct {
		name    string
		s       domain.StepName
		wantErr bool
	}{
		{
			name:    "validate a valid step name",
			s:       domain.StepNameIdNumber,
			wantErr: false,
		},
		{
			name:    "validate an invalid step name",
			s:       domain.StepName("invalid"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrUnsupportedStepName, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
