package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
)

func TestInstructionType_Target(t *testing.T) {
	tests := []struct {
		name string
		i    domain.InstructionType
		want domain.InstructionTarget
	}{
		{
			name: "target for DELIVERY_CUSTOMER",
			i:    domain.DeliveryCustomerInstructionType,
			want: domain.TargetCustomer,
		},
		{
			name: "target for DELIVERY_HOUSEHOLD",
			i:    domain.DeliveryHouseHoldInstructionType,
			want: domain.TargetHousehold,
		},
		{
			name: "target for DELIVERY_3RD_PARTY",
			i:    domain.Delivery3rdPartyInstructionType,
			want: domain.TargetThirdParty,
		},
		{
			name: "target for FAILED_DELIVERY_CUSTOMER",
			i:    domain.DeliveryCustomerFailedInstructionType,
			want: domain.TargetCustomer,
		},
		{
			name: "target for FAILED_DELIVERY_WH",
			i:    domain.DeliveryWHFailedInstructionType,
			want: domain.TargetWarehouse,
		},
		{
			name: "target for PICKUP_WH",
			i:    domain.PickupWHInstructionType,
			want: domain.TargetWarehouse,
		},
		{
			name: "target for FAILED_PICKUP_WH",
			i:    domain.PickupWHFailedInstructionType,
			want: domain.TargetWarehouse,
		},
		{
			name: "target for PICKUP_CUSTOMER",
			i:    domain.PickupCustomerInstructionType,
			want: domain.TargetCustomer,
		},
		{
			name: "target for FAILED_PICKUP_CUSTOMER",
			i:    domain.PickupCustomerFailedInstructionType,
			want: domain.TargetCustomer,
		},
		{
			name: "target for invalid type",
			i:    "INVALID_TYPE",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.i.Target()

			assert.Equal(t, tt.want, got)
		})
	}
}
