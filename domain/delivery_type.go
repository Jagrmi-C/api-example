package domain

import "errors"

// ErrUnsupportedDeliveryType is returned when delivery type is invalid.
var ErrUnsupportedDeliveryType = errors.New("unsupported delivery type")

type DeliveryType string

func (m DeliveryType) String() string {
	return string(m)
}

func (m DeliveryType) Validate() error {
	switch m {
	case DeliveryTypeDirect,
		DeliveryTypeReverse:
		return nil
	default:
		return ErrUnsupportedDeliveryType
	}
}

const (
	DeliveryTypeDirect  DeliveryType = "DIRECT"
	DeliveryTypeReverse DeliveryType = "REVERSE"
)
