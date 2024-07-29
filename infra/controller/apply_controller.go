package controller

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"gitlab.com/jc88/api-example/infra/errorx"
	"gitlab.com/jc88/api-example/usecase"
)

// ApplyRetailerConfigurationDI is struct for dependency injection.
type ApplyRetailerConfigurationDI struct {
	UC usecase.RetailerConfigurationsApplier
}

type applyRetailerConfiguration struct {
	uc usecase.RetailerConfigurationsApplier
}

// NewApplyRetailerConfiguration is a constructor.
func NewApplyRetailerConfiguration(di *ApplyRetailerConfigurationDI) *applyRetailerConfiguration {
	return &applyRetailerConfiguration{
		uc: di.UC,
	}
}

// Operation is a method to retrieve huma operation.
func (op *applyRetailerConfiguration) Operation() huma.Operation {
	return huma.Operation{
		OperationID:   "apply-retailer-configuration",
		Description:   "Apply up-to-date retailer configurations",
		Method:        http.MethodPost,
		Path:          "/config-manager/retailers/{retailerId}/configurations",
		Summary:       "Apply up-to-date retailer configurations",
		Tags:          []string{"Configurations"},
		DefaultStatus: http.StatusNoContent,
		Errors: []int{
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusInternalServerError,
		},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

// ApplyRetailerConfiguration is a huma handler.
func ApplyRetailerConfiguration(
	myUc *applyRetailerConfiguration,
) func(context.Context, *ApplyRetailerConfigurationReq) (*ApplyRetailerConfigurationRes, error) {
	return func(ctx context.Context, input *ApplyRetailerConfigurationReq) (*ApplyRetailerConfigurationRes, error) {
		labels := map[string]string{
			"retailerId": input.RetailerID.String(),
		}

		config, err := convertInputToDomainModel(input)
		if err != nil {
			return nil, renderErrResponse(err, nil, labels)
		}

		if err := myUc.uc.ApplyNewRetailerInstruction(ctx, config, input.Auth); err != nil {
			return nil, renderErrResponse(err, nil, labels)
		}

		response := ApplyRetailerConfigurationRes{
			ContentType:  "application/json; charset=utf-8",
			LastModified: time.Now().UTC(),
			Status:       http.StatusNoContent,
		}

		return &response, nil
	}
}

func renderErrResponse(err error, l *slog.Logger, labels map[string]string) error {
	switch errT := err.(type) {
	case *errorx.Error:
		switch errT.Code() {
		case errorx.ErrCodeRetailerNotExist:
			err = huma.Error404NotFound(errT.Error(), errT)
		case errorx.ErrCodeAuth0AuthMalformed:
			err = huma.Error401Unauthorized(errT.Error(), errT)
		case errorx.ErrorCodeInputValidation:
			err = huma.Error400BadRequest(errT.Error(), errT)
		case errorx.ErrCodeRetailerConfigurationAlreadyExist:
			err = huma.Error409Conflict(errT.Error(), errT)
		default:
			err = huma.Error500InternalServerError(errT.Error(), errT)
		}

	default:
		err = huma.Error500InternalServerError("unexpected error", err)
	}

	return err
}
