package domain

import "errors"

// ErrUnsupportedInstructionType is returned when the instruction type is invalid.
var ErrUnsupportedInstructionType = errors.New("unsupported instruction type")

// InstructionType is specific instruction type according to REMS domain.
type InstructionType string

// Instructions for delivery instructions.
const (
	DeliveryCustomerInstructionType       InstructionType = "DELIVERY_CUSTOMER"
	DeliveryHouseHoldInstructionType      InstructionType = "DELIVERY_HOUSEHOLD"
	Delivery3rdPartyInstructionType       InstructionType = "DELIVERY_3RD_PARTY"
	DeliveryCustomerFailedInstructionType InstructionType = "FAILED_DELIVERY_CUSTOMER"
	DeliveryWHFailedInstructionType       InstructionType = "FAILED_DELIVERY_WH"
)

// Instructions for pick up instructions.
const (
	PickupWHInstructionType             InstructionType = "PICKUP_WH"
	PickupWHFailedInstructionType       InstructionType = "FAILED_PICKUP_WH"
	PickupCustomerInstructionType       InstructionType = "PICKUP_CUSTOMER"
	PickupCustomerFailedInstructionType InstructionType = "FAILED_PICKUP_CUSTOMER"
)

func (i InstructionType) Validate() error {
	switch i {
	case DeliveryCustomerInstructionType,
		DeliveryHouseHoldInstructionType,
		Delivery3rdPartyInstructionType,
		DeliveryCustomerFailedInstructionType,
		DeliveryWHFailedInstructionType,
		PickupWHInstructionType,
		PickupWHFailedInstructionType,
		PickupCustomerInstructionType,
		PickupCustomerFailedInstructionType:
		return nil
	default:
		return ErrUnsupportedInstructionType
	}
}

func (i InstructionType) String() string {
	return string(i)
}

func (i InstructionType) Target() InstructionTarget {
	switch i {
	case DeliveryCustomerInstructionType:
		return TargetCustomer
	case DeliveryHouseHoldInstructionType:
		return TargetHousehold
	case Delivery3rdPartyInstructionType:
		return TargetThirdParty
	case DeliveryCustomerFailedInstructionType:
		return TargetCustomer
	case DeliveryWHFailedInstructionType:
		return TargetWarehouse
	case PickupWHInstructionType:
		return TargetWarehouse
	case PickupWHFailedInstructionType:
		return TargetWarehouse
	case PickupCustomerInstructionType:
		return TargetCustomer
	case PickupCustomerFailedInstructionType:
		return TargetCustomer
	default:
		return ""
	}
}
