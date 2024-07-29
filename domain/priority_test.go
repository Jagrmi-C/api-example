package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
)

func TestStepPriority_String(t *testing.T) {
	tests := []struct {
		name string
		s    domain.StepPriority
		want string
	}{
		{
			name: "string representation for ScanStep",
			s:    domain.ScanStep,
			want: "Scan",
		},
		{
			name: "string representation for CollectBarcodesStep",
			s:    domain.CollectBarcodesStep,
			want: "Collect Barcodes",
		},
		{
			name: "string representation for PhotoStep",
			s:    domain.PhotoStep,
			want: "Photo",
		},
		{
			name: "string representation for IDNumberStep",
			s:    domain.IDNumberStep,
			want: "ID Number",
		},
		{
			name: "string representation for SignatureStep",
			s:    domain.SignatureStep,
			want: "Signature",
		},
		{
			name: "string representation for PassCodeStep",
			s:    domain.PassCodeStep,
			want: "Passcode",
		},
		{
			name: "string representation for OTPStep",
			s:    domain.OTPStep,
			want: "OTP",
		},
		{
			name: "string representation for ManualConfirmationStep",
			s:    domain.ManualConfirmationStep,
			want: "Manual Confirmation",
		},
		{
			name: "string representation for ConfigurableTextStep",
			s:    domain.ConfigurableTextStep,
			want: "Configurable Text",
		},
		{
			name: "string representation for unknown StepPriority",
			s:    domain.StepPriority(100),
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.String()

			assert.Equal(t, tt.want, got)
		})
	}
}
