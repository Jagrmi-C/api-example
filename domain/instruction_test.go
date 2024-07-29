package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
)

func TestNewInstructionByTypeAndCountry(t *testing.T) {
	type args struct {
		countryStr         string
		instructionTypeStr string
	}
	tests := []struct {
		name        string
		args        args
		want        *domain.Instruction
		wantErr     bool
		expectedErr error
	}{
		{
			name: "success, existing combination instruction type+country  DELIVERY_CUSTOMER+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "DELIVERY_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country DELIVERY_CUSTOMER+FR",
			args: args{
				countryStr:         "FR",
				instructionTypeStr: "DELIVERY_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country DELIVERY_CUSTOMER+IT",
			args: args{
				countryStr:         "IT",
				instructionTypeStr: "DELIVERY_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country DELIVERY_CUSTOMER+PT",
			args: args{
				countryStr:         "PT",
				instructionTypeStr: "DELIVERY_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country DELIVERY_HOUSEHOLD+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "DELIVERY_HOUSEHOLD",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country FAILED_DELIVERY_CUSTOMER+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "FAILED_DELIVERY_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country DELIVERY_3RD_PARTY+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "DELIVERY_3RD_PARTY",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country FAILED_DELIVERY_WH+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "FAILED_DELIVERY_WH",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country PICKUP_WH+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "PICKUP_WH",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country PICKUP_CUSTOMER+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "PICKUP_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country PICKUP_CUSTOMER+FR",
			args: args{
				countryStr:         "FR",
				instructionTypeStr: "PICKUP_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country FAILED_PICKUP_CUSTOMER+FR",
			args: args{
				countryStr:         "FR",
				instructionTypeStr: "FAILED_PICKUP_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country FAILED_PICKUP_CUSTOMER+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "FAILED_PICKUP_CUSTOMER",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country FAILED_PICKUP_WH+ES",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "FAILED_PICKUP_WH",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country FAILED_PICKUP_WH+FR",
			args: args{
				countryStr:         "FR",
				instructionTypeStr: "FAILED_PICKUP_WH",
			},
			wantErr: false,
		},
		{
			name: "success, existing combination instruction type+country FAILED_PICKUP_WH+PT",
			args: args{
				countryStr:         "PT",
				instructionTypeStr: "FAILED_PICKUP_WH",
			},
			wantErr: false,
		},
		{
			name: "failed, invalid instruction type",
			args: args{
				countryStr:         "ES",
				instructionTypeStr: "PickupWHInstructionType",
			},
			wantErr:     true,
			expectedErr: domain.ErrUnsupportedInstructionType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewInstructionByTypeAndCountry(tt.args.countryStr, tt.args.instructionTypeStr)
			if tt.wantErr {
				assert.Error(t, err)

				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}

				return
			}

			assert.IsType(t, tt.want, got)
		})
	}
}

func TestInstructions_AddInstructionTypes(t *testing.T) {
	type args struct {
		country          domain.Country
		instructionTypes []domain.InstructionType
	}
	tests := []struct {
		name    string
		i       *domain.Instructions
		args    args
		wantErr bool
	}{
		{
			name: "success, instruction type with country were added",
			i: &domain.Instructions{
				domain.Instruction{
					Type:     domain.DeliveryCustomerInstructionType,
					Priority: 1,
					Target:   domain.TargetCustomer,
				},
			},
			args: args{
				country: domain.CountryES,
				instructionTypes: []domain.InstructionType{
					domain.DeliveryCustomerInstructionType,
					domain.DeliveryHouseHoldInstructionType,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.i.AddInstructionTypes(tt.args.country, tt.args.instructionTypes...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestNewInstruction(t *testing.T) {
	type args struct {
		instructionTypeStr string
		stepNames          []domain.StepName
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Instruction
		wantErr bool
	}{
		{
			name: "success, create a new instruction with multiple steps",
			args: args{
				instructionTypeStr: "DELIVERY_CUSTOMER",
				stepNames:          []domain.StepName{"Step1", "Step2", "Step3"},
			},
		},
		{
			name: "failed, invalid instruction type",
			args: args{
				instructionTypeStr: "failed",
				stepNames:          []domain.StepName{"Step1", "Step2", "Step3"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewInstruction(tt.args.instructionTypeStr, tt.args.stepNames)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.IsType(t, tt.want, got)
		})
	}
}
