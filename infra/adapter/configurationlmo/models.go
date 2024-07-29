package configurationlmo

import (
	"gitlab.com/jc88/api-example/domain"
)

type ConfigurationInstructionType string

const (
	PickupConfigurationInstructionType       ConfigurationInstructionType = "pickup"
	DeliveryConfigurationInstructionType     ConfigurationInstructionType = "delivery"
	FailPickupConfigurationInstructionType   ConfigurationInstructionType = "failPickup"
	FailDeliveryConfigurationInstructionType ConfigurationInstructionType = "failDelivery"
)

type RetailerConfigurationCreateReq struct {
	Instructions []Instruction `json:"instructions"`
}

func (req *RetailerConfigurationCreateReq) AddInstructionFromConfiguration(config *domain.Configuration) error {
	configurationInstructions := make([]Instruction, 0, len(config.Instructions))

	for i := range config.Instructions {
		instrReq := Instruction{
			Model:    ConvertToLMOModel(config.OperationalModel),
			Target:   string(config.Instructions[i].Target),
			Priority: config.Instructions[i].Priority,
		}

		// default instruction -> without country and service type
		country := string(config.Country)
		if country != "" {
			instrReq.Country = country
		}

		if config.ServiceType != "" {
			instrReq.ServiceType = config.ServiceType
		}

		instrReq.SetType(config.Instructions[i].Type)

		if err := instrReq.AddSteps(config.Instructions[i].Steps); err != nil {
			return err
		}

		configurationInstructions = append(configurationInstructions, instrReq)
	}

	if req.Instructions == nil {
		req.Instructions = configurationInstructions
	} else {
		req.Instructions = append(req.Instructions, configurationInstructions...)
	}

	return nil
}

type Instruction struct {
	Model          string                       `json:"model"`
	Type           ConfigurationInstructionType `json:"type"`
	Target         string                       `json:"target"`
	Priority       int8                         `json:"priority"`
	Steps          []Step                       `json:"steps"`
	Store          string                       `json:"store,omitempty"`
	Country        string                       `json:"country,omitempty"`
	ServiceType    string                       `json:"serviceType,omitempty"`
	Multiparcel    *bool                        `json:"multiparcel,omitempty"`
	Insured        *bool                        `json:"insured,omitempty"`
	CashOnDelivery *bool                        `json:"cashOnDelivery,omitempty"`
	Category       string                       `json:"category,omitempty"`
}

// SetModel sets the configuration model.
func ConvertToLMOModel(model domain.Model) string {
	switch model {
	case domain.ModelWarehouse:
		return "warehouse_model"
	case domain.ModelStore:
		return "store_model"
	default: // unknown
		return ""
	}
}

func (ins *Instruction) SetType(remsDomainType domain.InstructionType) {
	switch remsDomainType {
	case domain.DeliveryCustomerInstructionType,
		domain.DeliveryHouseHoldInstructionType,
		domain.Delivery3rdPartyInstructionType:
		ins.Type = DeliveryConfigurationInstructionType
	case domain.DeliveryCustomerFailedInstructionType, domain.DeliveryWHFailedInstructionType:
		ins.Type = FailDeliveryConfigurationInstructionType
	case domain.PickupWHInstructionType, domain.PickupCustomerInstructionType:
		ins.Type = PickupConfigurationInstructionType
	case domain.PickupWHFailedInstructionType, domain.PickupCustomerFailedInstructionType:
		ins.Type = FailPickupConfigurationInstructionType
	default: // unknown
	}
}

func (ins *Instruction) AddSteps(domainSteps []domain.InstructionStep) error {
	steps := make([]Step, 0, len(domainSteps))

	for i := range domainSteps {
		st := Step{
			Mandatory: domainSteps[i].Mandatory,
			Priority:  domainSteps[i].Priority,
			Attempts:  domainSteps[i].Attempts,
		}

		if err := st.SetLMOName(domainSteps[i].Name); err != nil {
			return err
		}

		if domainSteps[i].ExtraInfo != nil {
			st.ExtraInfo = make([]ExtraInfo, 0, len(domainSteps[i].ExtraInfo))

			for j := range domainSteps[i].ExtraInfo {
				st.ExtraInfo = append(st.ExtraInfo, ExtraInfo{
					Key:   domainSteps[i].ExtraInfo[j].Key,
					Value: domainSteps[i].ExtraInfo[j].Value,
				})
			}
		}

		if domainSteps[i].Fallback.Name != "" {
			st.Fallback = &Fallback{
				Name: domainSteps[i].Fallback.Name,
			}
		}

		steps = append(steps, st)
	}

	if len(steps) > 0 {
		ins.Steps = steps
	}

	return nil
}

type Step struct {
	Name      string      `json:"name"`
	Mandatory bool        `json:"mandatory"`
	Priority  int8        `json:"priority"`
	Attempts  int16       `json:"attempts"`
	ExtraInfo []ExtraInfo `json:"extraInfo,omitempty"`
	Fallback  *Fallback   `json:"fallback,omitempty"`
}

type Fallback struct {
	Name string `json:"name,omitempty"`
}

func (s *Step) SetLMOName(step domain.StepName) error {
	switch step {
	case domain.StepNameScan:
		s.Name = "scan"
	case domain.StepNameCollectBarcode:
		s.Name = "collectBarcodes"
	case domain.StepNamePhoto:
		s.Name = "photo"
	case domain.StepNameIdNumber:
		s.Name = "identityDocumentNumber"
	case domain.StepNameSignature:
		s.Name = "signature"
	case domain.StepNamePassCode:
		s.Name = "retailerPasscode"
	case domain.StepNameOTP:
		s.Name = "otp"
	case domain.StepNameManualConfirmation:
		s.Name = "manualConfirmation"
	case domain.StepNameConfigurableText:
		s.Name = string(step)
	default:
		return domain.ErrUnsupportedStepName
	}

	return nil
}

type ExtraInfo struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func NewRetailerConfigurationCreateReq(config *domain.RetailerConfiguration) (*RetailerConfigurationCreateReq, error) {
	req := RetailerConfigurationCreateReq{
		Instructions: make([]Instruction, 0, len(config.Configurations)),
	}

	// according to documentation -> required fields
	// - model
	// - type
	// - target
	// - priority
	// - steps

	for _, instr := range config.Configurations {
		if err := req.AddInstructionFromConfiguration(&instr); err != nil {
			return nil, err
		}
	}

	return &req, nil
}

// GetConfigurationsByRetailerRespo describes a response from RIS system.
type GetConfigurationsByRetailerRespo struct {
	Data  GetConfigurationsByRetailerRespoData `json:"data"`
	Error *ResponseError                       `json:"error"`
}

type GetConfigurationsByRetailerRespoData struct {
	RetailerID   string                                           `json:"id"`
	Instructions GetConfigurationsByRetailerRespoDataInstructions `json:"instructions"`
}

type GetConfigurationsByRetailerRespoDataInstructions struct {
	ID       string `json:"id"`
	Model    string `json:"model"`
	Type     string `json:"type"`
	Target   string `json:"target"`
	Priority int16  `json:"priority"`
	Attempts int16  `json:"attempts"`
	Steps    []struct {
		Name      string `json:"name"`
		Mandatory bool   `json:"mandatory"`
		Priority  int16  `json:"priority"`
		ExtraInfo []any  `json:"extraInfo"`
		Fallback  struct {
			Name string `json:"name"`
		} `json:"fallback"`
	} `json:"steps"`
	OrderType string `json:"order_type"`
	Country   string `json:"country,omitempty"`
}

type RetrieveProcessByIDDetailsRespoData struct {
	ID         string                                        `json:"id"`
	Attributes RetrieveProcessByIDDetailsRespoDataAttributes `json:"attributes"`
}

type RetrieveProcessByIDDetailsRespoDataAttributes struct {
	FilePattern string  `json:"file_pattern"`
	Status      string  `json:"status"`
	Timestamp   float64 `json:"timestamp"`
}

// ResponseError describes an error response.
type ResponseError struct {
	Msg  string `json:"msg,omitempty"`
	Code string `json:"code,omitempty"`
}
