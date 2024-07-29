package domain

import (
	uuid "github.com/satori/go.uuid"
)

// Instruction is a entity representing instructions for driver.
type Instruction struct {
	ID       uuid.UUID
	Type     InstructionType // by type is it possible to calculate target and country i
	Priority int8
	Target   InstructionTarget
	Steps    []InstructionStep
}

// NewInstructionByTypeAndCountry generates a new Instruction by type and country.
// Contains some business logic to determine how the instruction can be generated
// by country and REMS instruction type.
func NewInstructionByTypeAndCountry(countryStr, instructionTypeStr string) (*Instruction, error) {
	country := Country(countryStr)
	if err := country.Validate(); err != nil {
		return nil, err
	}

	instructionType := InstructionType(instructionTypeStr)
	if err := instructionType.Validate(); err != nil {
		return nil, err
	}

	instruction := Instruction{
		Type: instructionType,
	}

	if _, err := instruction.PrepareInstructionBYTargetAndCountryFromSpreadsheet(instructionType, country); err != nil {
		return nil, err
	}

	return &instruction, nil
}

// NewInstruction is constructor.
func NewInstruction(instructionTypeStr string, stepNames []StepName) (*Instruction, error) {
	instructionType := InstructionType(instructionTypeStr)
	if err := instructionType.Validate(); err != nil {
		return nil, err
	}

	i := Instruction{}
	i.PrepareInstructionByType(instructionType)
	i.PrepareSteps(stepNames...)

	return &i, nil
}

// PrepareSteps is a helper method to prepare steps by step names.
func (i *Instruction) PrepareSteps(stepNames ...StepName) {
	steps := make([]InstructionStep, 0, len(stepNames))

	for i := range stepNames {
		steps = append(steps, InstructionStep{
			Name:      stepNames[i],
			Mandatory: true,
			Attempts:  1,
			Priority:  stepNames[i].Priority(),
		})
	}

	i.Steps = steps
}

// SetTargetWithPriority is a helper method to set target and priority for a given instruction target.
func (i *Instruction) SetTargetWithPriority(target InstructionTarget) {
	i.Target = target
	i.Priority = target.Priority()
}

func (i *Instruction) PrepareInstructionByType(iType InstructionType) {
	i.Type = iType
	i.Target = iType.Target()
	i.Priority = iType.Target().Priority()

	i.PrepareSteps(StepNameIdNumber, StepNameSignature)
}

// PrepareInstructionBYTargetAndCountryFromSpreadsheet converts rems instruction type to LMO instruction type,
// that should knows about country mandatory steps.
// nolint:all
func (i *Instruction) PrepareInstructionBYTargetAndCountryFromSpreadsheet(iType InstructionType, country Country) (DeliveryType, error) {
	var deliveryType DeliveryType

	i.Type = iType

	switch i.Type {
	case DeliveryCustomerInstructionType:
		i.SetTargetWithPriority(TargetCustomer)
		deliveryType = DeliveryTypeDirect

		switch country {
		case CountryES:
			i.PrepareSteps(StepNameIdNumber, StepNameSignature)
		case CountryFR:
			i.PrepareSteps(StepNameSignature, StepNameSignature)
		case CountryIT:
			i.PrepareSteps(StepNameSignature, StepNamePhoto)
		case CountryPT:
			i.PrepareSteps(StepNameIdNumber, StepNameSignature)
		}
	case DeliveryCustomerFailedInstructionType:
		i.SetTargetWithPriority(TargetCustomer)
		deliveryType = DeliveryTypeDirect

		switch country {
		case CountryES:
			i.PrepareSteps(StepNamePhoto)
		case CountryFR:
			i.PrepareSteps(StepNamePhoto)
		case CountryIT:
			i.PrepareSteps(StepNamePhoto)
		case CountryPT:
			i.PrepareSteps(StepNamePhoto)
		}
	case DeliveryHouseHoldInstructionType:
		i.SetTargetWithPriority(TargetHousehold)
		deliveryType = DeliveryTypeDirect

		switch country {
		case CountryES:
			i.PrepareSteps(StepNameIdNumber, StepNameSignature)
		case CountryFR:
			i.PrepareSteps(StepNameSignature, StepNamePhoto)
		case CountryIT:
			i.PrepareSteps(StepNameSignature, StepNamePhoto)
		case CountryPT:
			i.PrepareSteps(StepNameIdNumber, StepNameSignature)
		}
	case Delivery3rdPartyInstructionType:
		i.SetTargetWithPriority(TargetThirdParty)
		deliveryType = DeliveryTypeDirect

		switch country {
		case CountryES:
			i.PrepareSteps(StepNameIdNumber, StepNameSignature)
		case CountryFR:
			i.PrepareSteps(StepNameSignature, StepNamePhoto)
		case CountryIT:
			i.PrepareSteps(StepNameSignature, StepNamePhoto)
		case CountryPT:
			i.PrepareSteps(StepNameIdNumber, StepNameSignature)
		}
	case DeliveryWHFailedInstructionType:
		i.SetTargetWithPriority(TargetWarehouse)
		deliveryType = DeliveryTypeReverse

		switch country {
		case CountryES:
			i.PrepareSteps(StepNamePhoto)
		case CountryFR:
		case CountryIT:
		case CountryPT:
			i.PrepareSteps(StepNamePhoto)
		}
	case PickupWHInstructionType:
		i.SetTargetWithPriority(TargetWarehouse)
		deliveryType = DeliveryTypeDirect

		switch country {
		case CountryES:
			i.PrepareSteps(StepNameScan)
		case CountryFR:
			i.PrepareSteps(StepNameScan)
		case CountryIT:
			i.PrepareSteps(StepNameScan)
		case CountryPT:
			i.PrepareSteps(StepNameScan)
		}
	case PickupWHFailedInstructionType:
		i.SetTargetWithPriority(TargetWarehouse)
		deliveryType = DeliveryTypeDirect

		switch country {
		case CountryES:
			i.PrepareSteps(StepNamePhoto)
		case CountryFR:
			i.PrepareSteps(StepNamePhoto)
		case CountryIT:
			i.PrepareSteps(StepNamePhoto)
		case CountryPT:
			i.PrepareSteps(StepNamePhoto)
		}
	case PickupCustomerInstructionType:
		i.SetTargetWithPriority(TargetCustomer)
		deliveryType = DeliveryTypeReverse

		switch country {
		case CountryES:
			i.PrepareSteps(StepNameScan, StepNamePhoto)
		case CountryFR:
		case CountryIT:
		case CountryPT:
		}
	case PickupCustomerFailedInstructionType:
		i.SetTargetWithPriority(TargetCustomer)
		deliveryType = DeliveryTypeReverse

		switch country {
		case CountryES:
			i.PrepareSteps(StepNamePhoto)
		case CountryFR:
			i.PrepareSteps(StepNamePhoto)
		case CountryIT:
			i.PrepareSteps(StepNamePhoto)
		case CountryPT:
			i.PrepareSteps(StepNamePhoto)
		}
	default:
		return deliveryType, ErrUnsupportedInstructionType
	}

	return deliveryType, nil
}

// Instructions is a custom type that can be used to identify arrays of instructions.
type Instructions []Instruction

// AddInstructionTypes is a helper method to add instruction/s by instruction type and country.
func (i *Instructions) AddInstructionTypes(country Country, instructionTypes ...InstructionType) error {
	for _, iType := range instructionTypes {
		for _, inst := range *i {
			if inst.Type == iType {
				continue
			}
		}

		instruction, err := NewInstructionByTypeAndCountry(country.String(), iType.String())
		if err != nil {
			return err
		}

		*i = append(*i, *instruction)
	}

	return nil
}

// InstructionStep is value object representing instructions step.
type InstructionStep struct {
	Name      StepName
	Mandatory bool
	Attempts  int16
	Priority  int8
	ExtraInfo []ExtraInfo
	Fallback  FallbackStep
}

// FallbackStep is value object representing instructions step with fallback.
type FallbackStep struct {
	Name string
}

// ExtraInfo is a value object representing additional extra info.
type ExtraInfo struct {
	Key   string
	Value string
}
