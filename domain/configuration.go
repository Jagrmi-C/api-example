package domain

import uuid "github.com/satori/go.uuid"

// RetailerConfiguration is an aggregator that includes information about retailers configuration,
// and including instruction.
type RetailerConfiguration struct {
	RetailerID     uuid.UUID
	Configurations []Configuration
}

type Configuration struct {
	RetailerID       uuid.UUID
	OperationalModel Model
	ServiceType      string
	DeliveryType     DeliveryType
	Country          Country
	Instructions     Instructions
}

// NewDefaultPickUPWarehouse create a new default configuration for REMS type PICKUP_WH.
func NewDefaultPickUPWarehouse(retailerID uuid.UUID) *Configuration {
	step := InstructionStep{
		Name:      StepNameScan,
		Mandatory: true,
		Attempts:  1,
		Priority:  StepNameScan.Priority(),
		Fallback: FallbackStep{
			Name: "manualConfirmation",
		},
	}

	instr := Instruction{
		Type:     PickupWHInstructionType,
		Target:   TargetWarehouse,
		Priority: TargetWarehouse.Priority(),
		Steps:    []InstructionStep{step},
	}

	config := Configuration{
		RetailerID:       retailerID,
		OperationalModel: ModelWarehouse,
		DeliveryType:     DeliveryTypeDirect,
		Instructions:     []Instruction{instr},
	}

	return &config
}

// NewDefaultPickUPWarehouse create a new default configuration for REMS type FAILED_PICKUP_WH.
func NewDefaultFailedPickUPWarehouse(retailerID uuid.UUID) *Configuration {
	instr := Instruction{
		Type:     PickupWHFailedInstructionType,
		Target:   TargetWarehouse,
		Priority: TargetWarehouse.Priority(),
	}

	instr.PrepareSteps(StepNamePhoto)

	config := Configuration{
		RetailerID:       retailerID,
		OperationalModel: ModelWarehouse,
		DeliveryType:     DeliveryTypeReverse,
		Instructions:     []Instruction{instr},
	}

	return &config
}

// NewDefaultFailedDeliveryWarehouse create a new default configuration for REMS type FAILED_DELIVERY_WH.
func NewDefaultFailedDeliveryWarehouse(retailerID uuid.UUID) *Configuration {
	instr := Instruction{
		Type:     DeliveryWHFailedInstructionType,
		Target:   TargetWarehouse,
		Priority: TargetWarehouse.Priority(),
	}

	instr.PrepareSteps(StepNamePhoto)

	config := Configuration{
		RetailerID:       retailerID,
		OperationalModel: ModelWarehouse,
		DeliveryType:     DeliveryTypeDirect,
		Instructions:     []Instruction{instr},
	}

	return &config
}

// AddDefaultConfigurations is a method to add default mandatory configurations.
func (c *RetailerConfiguration) AddDefaultConfigurations() {
	config1 := NewDefaultPickUPWarehouse(c.RetailerID)
	config2 := NewDefaultFailedPickUPWarehouse(c.RetailerID)
	config3 := NewDefaultFailedDeliveryWarehouse(c.RetailerID)

	c.Configurations = append(c.Configurations, *config1, *config2, *config3)
}
