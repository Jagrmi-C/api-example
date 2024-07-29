package domain

import "errors"

// ErrUnsupportedInstructionType is returned when the instruction type is invalid.
var ErrUnsupportedStepName = errors.New("unsupported step type according to name")

// StepName is custom step name type.
type StepName string

func (s StepName) Validate() error {
	switch s {
	case StepNameScan,
		StepNameCollectBarcode,
		StepNamePhoto,
		StepNameIdNumber,
		StepNameSignature,
		StepNamePassCode,
		StepNameOTP,
		StepNameManualConfirmation,
		StepNameConfigurableText:
		return nil
	default:
		return ErrUnsupportedStepName
	}
}

func (s StepName) String() string {
	return string(s)
}

// Priority is custom method that determines priority of the step.
// Scan > Collect Barcodes > Photo > ID Number > Signature > Passcode > OTP  > Manual Confirmation > Configurable Text.
func (s StepName) Priority() int8 {
	switch s {
	case StepNameScan:
		return 1
	case StepNameCollectBarcode:
		return 2
	case StepNamePhoto:
		return 3
	case StepNameIdNumber:
		return 4
	case StepNameSignature:
		return 5
	case StepNamePassCode:
		return 6
	case StepNameOTP:
		return 7
	case StepNameManualConfirmation:
		return 8
	case StepNameConfigurableText:
		return 9
	default:
		return 0
	}
}

const (
	StepNameScan               StepName = "SCAN"
	StepNameCollectBarcode     StepName = "COLLECT_BARCODES"
	StepNamePhoto              StepName = "PHOTO"
	StepNameIdNumber           StepName = "ID_NUMBER"
	StepNameSignature          StepName = "SIGNATURE"
	StepNamePassCode           StepName = "PASSCODE"
	StepNameOTP                StepName = "OTP"
	StepNameManualConfirmation StepName = "MANUAL_CONFIRMATION"
	StepNameConfigurableText   StepName = "CONFIGURABLE_TEXT"
)
