package usecase_test

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jc88/api-example/domain"
	"gitlab.com/jc88/api-example/usecase"
)

type mocker struct {
	m *MockRetailerInstructionsManager
}

func NewMocker(repo *MockRetailerInstructionsManager) *mocker {
	return &mocker{
		m: repo,
	}
}

func Test_retailerConfigurationsApplier_ApplyRetailerConfigurations(t *testing.T) {
	retailerID, err := uuid.NewV4()
	assert.NoError(t, err)

	auth0Token := "Bearer tets"

	retailerConfig := domain.RetailerConfiguration{
		RetailerID: retailerID,
		Configurations: []domain.Configuration{
			{
				RetailerID:   retailerID,
				ServiceType:  "ST1",
				DeliveryType: domain.DeliveryTypeDirect,
				Country:      domain.CountryES,
				Instructions: domain.Instructions{
					{},
				},
			},
		},
	}

	type args struct {
		config     *domain.RetailerConfiguration
		auth0Token string
	}

	tests := []struct {
		name         string
		prepareMocks func(ctx context.Context, factory *mocker)
		args         args
		wantErr      bool
	}{
		{
			name: "success, applying new state of retailer configurations",
			args: args{
				config:     &retailerConfig,
				auth0Token: auth0Token,
			},
			prepareMocks: func(ctx context.Context, factory *mocker) {
				factory.m.EXPECT().
					DeleteConfigurationsByRetailerID(ctx, retailerID, auth0Token).
					Once().
					Return(nil)

				factory.m.EXPECT().
					Create(ctx, &retailerConfig, auth0Token).
					Once().
					Return(nil)
			},
		},
		{
			name: "failed, Create method returns an unexpected error",
			args: args{
				config:     &retailerConfig,
				auth0Token: auth0Token,
			},
			prepareMocks: func(ctx context.Context, factory *mocker) {
				factory.m.EXPECT().
					DeleteConfigurationsByRetailerID(ctx, retailerID, auth0Token).
					Once().
					Return(nil)

				factory.m.EXPECT().
					Create(ctx, &retailerConfig, auth0Token).
					Maybe().
					Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed, DeleteConfigurationsByRetailerID method returns an unexpected error",
			args: args{
				config: &domain.RetailerConfiguration{
					RetailerID: retailerID,
					Configurations: []domain.Configuration{
						{
							RetailerID:   retailerID,
							ServiceType:  "ST1",
							DeliveryType: domain.DeliveryTypeDirect,
							Country:      domain.CountryES,
							Instructions: domain.Instructions{
								{},
							},
						},
					},
				},
				auth0Token: auth0Token,
			},
			prepareMocks: func(ctx context.Context, factory *mocker) {
				factory.m.EXPECT().
					DeleteConfigurationsByRetailerID(ctx, retailerID, auth0Token).
					Once().
					Return(assert.AnError)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockedAdapter := NewMockRetailerInstructionsManager(t)

			factory := NewMocker(mockedAdapter)
			assert.NotNil(t, tt.prepareMocks)
			tt.prepareMocks(ctx, factory)

			m := usecase.NewRetailerConfigurationsApplier(&usecase.RetailerConfigurationsApplierDI{
				LMOAdapter: factory.m,
			})

			err := m.ApplyNewRetailerInstruction(ctx, tt.args.config, tt.args.auth0Token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
