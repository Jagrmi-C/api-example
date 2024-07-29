package domain

import "errors"

// ErrUnsupportedRetailerModel is returned when model is invalid.
var ErrUnsupportedRetailerModel = errors.New("unsupported retailer model")

type Model string

func (m Model) String() string {
	return string(m)
}

func (m Model) Validate() error {
	switch m {
	case ModelWarehouse, ModelStore:
		return nil
	default:
		return ErrUnsupportedRetailerModel
	}
}

const (
	ModelWarehouse Model = "WAREHOUSE"
	ModelStore     Model = "STORE"
)
