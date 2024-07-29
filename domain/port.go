package domain

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

// RetailerInstructionsManager represents communication that allows to change
// retailer instructions in LMO system.
type RetailerInstructionsManager interface {
	// Create is a method to create new configurations in LMO for the retailer.
	Create(ctx context.Context, instructions *RetailerConfiguration, auth0token string) error
	// DeleteConfigurationsByRetailerID is a method to delete all configurations for the retailer.
	DeleteConfigurationsByRetailerID(ctx context.Context, retailerID uuid.UUID, auth0token string) error
}
