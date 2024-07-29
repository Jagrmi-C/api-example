package domain

import "errors"

// ErrUnsupportedInstructionType is returned when the instruction type is invalid.
var ErrUnsupportedInstructionTarget = errors.New("unsupported instruction target")

type InstructionTarget string

const (
	TargetCustomer   InstructionTarget = "customer"
	TargetWarehouse  InstructionTarget = "warehouse"
	TargetStore      InstructionTarget = "store"
	TargetThirdParty InstructionTarget = "thirdParty"
	TargetSafePlace  InstructionTarget = "safePlace"
	TargetHousehold  InstructionTarget = "household"
)

// Customer, Third Party, Safe Place, Household, Warehouse

func (c InstructionTarget) IsActionTargetValid() error {
	switch c {
	case TargetCustomer,
		TargetWarehouse,
		TargetStore,
		TargetSafePlace,
		TargetThirdParty,
		TargetHousehold:
		return nil
	default:
		return ErrUnsupportedInstructionTarget
	}
}

func (c InstructionTarget) Priority() int8 {
	switch c {
	case TargetCustomer:
		return 1
	case TargetThirdParty:
		return 2
	case TargetSafePlace:
		return 3
	case TargetHousehold:
		return 4
	case TargetWarehouse:
		return 5
	case TargetStore:
		return 6 // clarify this target
	default:
		return 0
	}
}
