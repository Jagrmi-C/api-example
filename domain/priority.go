package domain

// StepPriority is the priority of the step as business rule.
// Scan > Collect Barcodes > Photo > ID Number > Signature > Passcode > OTP > Manual Confirmation > Configurable Text.
type StepPriority int

const (
	ScanStep StepPriority = iota + 1
	CollectBarcodesStep
	PhotoStep
	IDNumberStep
	SignatureStep
	PassCodeStep
	OTPStep
	ManualConfirmationStep
	ConfigurableTextStep
)

// String is a string representation of a StepPriority.
func (s StepPriority) String() string {
	switch s {
	case ScanStep:
		return "Scan"
	case CollectBarcodesStep:
		return "Collect Barcodes"
	case PhotoStep:
		return "Photo"
	case IDNumberStep:
		return "ID Number"
	case SignatureStep:
		return "Signature"
	case PassCodeStep:
		return "Passcode"
	case OTPStep:
		return "OTP"
	case ManualConfirmationStep:
		return "Manual Confirmation"
	case ConfigurableTextStep:
		return "Configurable Text"
	default:
		return "unknown"
	}
}
