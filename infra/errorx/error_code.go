package errorx

// ErrorCode defines supported error codes.
type ErrorCode string

const (
	ErrCodeInternalError                     ErrorCode = "LMO-ERR-0001"
	ErrCodeAuth0AuthMalformed                ErrorCode = "LMO-ERR-0002"
	ErrCodeRetailerNotExist                  ErrorCode = "LMO-ERR-0011"
	ErrCodeRetailerConfigurationNotExist     ErrorCode = "LMO-ERR-0013"
	ErrCodeRetailerConfigurationAlreadyExist ErrorCode = "LMO-ERR-0014"
	ErrCodeCountryNotExist                   ErrorCode = "LMO-ERR-0053"

	ErrorCodeInputValidation ErrorCode = "LMO-ERR-0105"
)

// String returns a string representation of the error code.
func (e ErrorCode) String() string {
	return string(e)
}

// ErrMsg returns a string representation of the error message.
func (e ErrorCode) ErrMsg() string {
	switch e {
	case ErrCodeInternalError:
		return "Internal server error"
	case ErrCodeAuth0AuthMalformed:
		return "Unable to make authentication. Please try again later"
	case ErrCodeRetailerNotExist:
		return "The retailer doesn't exists"
	case ErrCodeCountryNotExist:
		return "The country doesn't exists"
	case ErrCodeRetailerConfigurationNotExist:
		return "Retailer configuration doesn't exist"
	case ErrCodeRetailerConfigurationAlreadyExist:
		return "Retailer configuration already exists"
	default:
		return "Unknown Error"
	}
}
