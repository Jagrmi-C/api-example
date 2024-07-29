// nolint:all
package controller

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"gitlab.com/jc88/api-example/domain"
)

// ApplyRetailerConfigurationReq represents the review operation request.
type ApplyRetailerConfigurationReq struct {
	RetailerID uuid.UUID `path:"retailerId" format:"uuid"`
	Auth       string    `header:"X-Forwarded-Authorization" required:"true"`
	Body       struct {
		Data []Data `json:"data" required:"true"`
	}
}

type Data struct {
	Type       string          `json:"type" required:"true" default:"configurations"`
	Attributes ApplyAttributes `json:"attributes" required:"true"`
}

// ApplyAttributes is a helper struct that contains the attributes and skip error for extra attributes.
type ApplyAttributes struct {
	OperationalModel string        `json:"operationalModel" required:"true" enum:"WAREHOUSE,STORE" doc:"Operational Model"`
	ServiceType      string        `json:"serviceType" required:"true" doc:"ServiceType"`
	DeliveryType     string        `json:"deliveryType" required:"true" enum:"DIRECT,REVERSE" doc:"DeliveryType"`
	Country          string        `json:"country" required:"true" enum:"ES,IT,PT,FR" doc:"Country"`
	Instructions     []Instruction `json:"instructions" required:"true" doc:"Instructions"`
	_                struct{}      `json:"-" additionalProperties:"true"`
}

// Instruction is helper struct for req.
type Instruction struct {
	Priority int16  `json:"priority" required:"false" minimum:"0" maximum:"20" doc:"Priority"`
	Attempts int16  `json:"attempts"  required:"false" minimum:"0" doc:"Attempts"`
	Type     string `json:"type" required:"true" enum:"PICKUP_WH,FAILED_PICKUP_WH,DELIVERY_CUSTOMER,DELIVERY_HOUSEHOLD,DELIVERY_3RD_PARTY,FAILED_DELIVERY_CUSTOMER,PICKUP_CUSTOMER,FAILED_PICKUP_CUSTOMER,FAILED_DELIVERY_WH" doc:"Instruction Type"`
	Steps    []Step `json:"steps" required:"true" doc:"Steps"`
}

// Step is a helper struct for req.
type Step struct {
	Name string `json:"name" required:"true" enum:"SCAN,COLLECT_BARCODES,PHOTO,ID_NUMBER,SIGNATURE,PASSCODE,OTP,MANUAL_CONFIRMATION,CONFIGURABLE_TEXT" doc:"Name of the step"`
}

// ApplyRetailerConfigurationRes is response struct.
type ApplyRetailerConfigurationRes struct {
	ContentType  string    `header:"Content-Type"`
	LastModified time.Time `header:"Last-Modified"`
	Status       int
}

func convertInputToDomainModel(input *ApplyRetailerConfigurationReq) (*domain.RetailerConfiguration, error) {
	retailerConfigurations := make([]domain.Configuration, 0, len(input.Body.Data))

	for i := range input.Body.Data {
		retailerConfig := domain.Configuration{
			RetailerID:       input.RetailerID,
			OperationalModel: convertDeliveryModel(input.Body.Data[i].Attributes.OperationalModel),
			Country:          domain.Country(input.Body.Data[i].Attributes.Country),
			ServiceType:      input.Body.Data[i].Attributes.ServiceType,
			DeliveryType:     convertDeliveryType(input.Body.Data[i].Attributes.DeliveryType),
		}

		instructions := make([]domain.Instruction, 0, len(input.Body.Data[i].Attributes.Instructions))

		for _, instr := range input.Body.Data[i].Attributes.Instructions {
			instruction, err := domain.NewInstruction(
				instr.Type,
				convertStepName(instr.Steps),
			)
			if err != nil {
				return nil, err
			}

			instructions = append(instructions, *instruction)
		}

		retailerConfig.Instructions = instructions

		retailerConfigurations = append(retailerConfigurations, retailerConfig)
	}

	retailerConfig := domain.RetailerConfiguration{
		RetailerID:     input.RetailerID,
		Configurations: retailerConfigurations,
	}

	return &retailerConfig, nil
}

func convertDeliveryType(deliveryType string) domain.DeliveryType {
	switch deliveryType {
	case "DIRECT":
		return domain.DeliveryTypeDirect
	case "REVERSE":
		return domain.DeliveryTypeReverse
	default:
		return domain.DeliveryType("")
	}
}

func convertDeliveryModel(deliveryModel string) domain.Model {
	switch deliveryModel {
	case "WAREHOUSE":
		return domain.ModelWarehouse
	case "STORE":
		return domain.ModelStore
	default:
		return domain.Model("")
	}
}

func convertStepName(steps []Step) []domain.StepName {
	res := make([]domain.StepName, 0, len(steps))

	for i := range steps {
		res = append(res, domain.StepName(steps[i].Name))
	}

	return res
}
