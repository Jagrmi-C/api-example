package domain

import "errors"

// ErrUnsupportedDeliveryType is returned when delivery type is invalid.
var ErrUnsupportedCountry = errors.New("unsupported country")

type Country string

func (m Country) String() string {
	return string(m)
}

func (m Country) Validate() error {
	switch m {
	case CountryES, CountryIT, CountryFR, CountryPT:
		return nil
	default:
		return ErrUnsupportedCountry
	}
}

const (
	CountryES Country = "ES"
	CountryIT Country = "IT"
	CountryFR Country = "FR"
	CountryPT Country = "PT"
)
