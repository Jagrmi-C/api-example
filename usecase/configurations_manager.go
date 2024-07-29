package usecase

import (
	"context"
	"log/slog"

	"gitlab.com/jc88/api-example/domain"
)

// RetailerConfigurationsApplier is an interface that provides functionality to apply REMS configurations.
type RetailerConfigurationsApplier interface {
	// ApplyNewRetailerInstruction is a method that apply REMS configurations in the system.
	ApplyNewRetailerInstruction(ctx context.Context, config *domain.RetailerConfiguration, auth0token string) error
}

type (
	RetailerConfigurationsApplierDI struct {
		Logger     *slog.Logger
		LMOAdapter domain.RetailerInstructionsManager
	}

	retailerConfigurationsApplier struct {
		logger  *slog.Logger
		adapter domain.RetailerInstructionsManager
	}
)

// NewRetailerConfigurationsApplier is constructor.
func NewRetailerConfigurationsApplier(di *RetailerConfigurationsApplierDI) *retailerConfigurationsApplier {
	return &retailerConfigurationsApplier{
		logger:  di.Logger,
		adapter: di.LMOAdapter,
	}
}

func (m *retailerConfigurationsApplier) ApplyNewRetailerInstruction(ctx context.Context, config *domain.RetailerConfiguration, auth0token string) error {
	config.AddDefaultConfigurations()

	if err := m.adapter.DeleteConfigurationsByRetailerID(ctx, config.RetailerID, auth0token); err != nil {
		return err
	}

	if err := m.adapter.Create(ctx, config, auth0token); err != nil {
		return err
	}

	return nil
}
