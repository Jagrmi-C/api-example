package controller_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/jc88/api-example/infra/controller"
	"gitlab.com/jc88/api-example/infra/errorx"
)

func prepareApplyRetailerConfigurationRequest(filePath string) (io.ReadCloser, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	reader := io.NopCloser(bytes.NewReader(data))

	return reader, nil
}

type mockedFactory struct {
	uc *MockConfigurationApplier
}

func NewMockedFactory(
	uc *MockConfigurationApplier,
) *mockedFactory {
	return &mockedFactory{
		uc: uc,
	}
}

func TestApplyRetailerConfiguration(t *testing.T) {
	retailerID, err := uuid.NewV4()
	assert.NoError(t, err)

	someURL := "http://localhost"

	type args struct {
		pathParam string
		filePath  string
	}
	tests := []struct {
		name               string
		args               args
		prepareMocks       func(factory *mockedFactory)
		expectedStatusCode int
	}{
		{
			name: "success path",
			args: args{
				pathParam: retailerID.String(),
				filePath:  "../../testdata/stubs/configurations/req/create_retailer_config_200.json",
			},
			prepareMocks: func(factory *mockedFactory) {
				factory.uc.EXPECT().
					ApplyNewRetailerInstruction(mock.Anything, mock.Anything, mock.Anything).
					Once().
					Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "failed, input validation errors",
			args: args{
				pathParam: retailerID.String(),
				filePath:  "../../testdata/stubs/configurations/req/create_retailer_config_422_invalid_struct.json",
			},
			prepareMocks:       func(factory *mockedFactory) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "failed, invalid retailer id in path",
			args: args{
				pathParam: "invalid",
				filePath:  "../../testdata/stubs/configurations/req/create_retailer_config_200.json",
			},
			prepareMocks:       func(factory *mockedFactory) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "failed, authorized request to downstream",
			args: args{
				pathParam: retailerID.String(),
				filePath:  "../../testdata/stubs/configurations/req/create_retailer_config_200.json",
			},
			prepareMocks: func(factory *mockedFactory) {
				factory.uc.EXPECT().
					ApplyNewRetailerInstruction(mock.Anything, mock.Anything, mock.Anything).
					Once().
					Return(errorx.NewErrorf(
						errorx.ErrCodeAuth0AuthMalformed,
						fmt.Sprintf("unauthorizedAccessToMsg: %s", someURL),
					))
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "failed, business domain validation error",
			args: args{
				pathParam: retailerID.String(),
				filePath:  "../../testdata/stubs/configurations/req/create_retailer_config_200.json",
			},
			prepareMocks: func(factory *mockedFactory) {
				factory.uc.EXPECT().
					ApplyNewRetailerInstruction(mock.Anything, mock.Anything, mock.Anything).
					Once().
					Return(errorx.NewErrorf(
						errorx.ErrorCodeInputValidation,
						"some fields cannot be more that 200 characters",
					))
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "failed, with unexpected error",
			args: args{
				pathParam: retailerID.String(),
				filePath:  "../../testdata/stubs/configurations/req/create_retailer_config_200.json",
			},
			prepareMocks: func(factory *mockedFactory) {
				factory.uc.EXPECT().
					ApplyNewRetailerInstruction(mock.Anything, mock.Anything, mock.Anything).
					Once().
					Return(assert.AnError)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "failed, retailer doesn't exist error",
			args: args{
				pathParam: retailerID.String(),
				filePath:  "../../testdata/stubs/configurations/req/create_retailer_config_200.json",
			},
			prepareMocks: func(factory *mockedFactory) {
				factory.uc.EXPECT().
					ApplyNewRetailerInstruction(mock.Anything, mock.Anything, mock.Anything).
					Once().
					Return(errorx.NewErrorf(
						errorx.ErrCodeRetailerNotExist,
						"retailer doesn't exist",
					))
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, api := humatest.New(t)

			uc := NewMockConfigurationApplier(t)

			factory := NewMockedFactory(uc)

			assert.NotNil(t, tt.prepareMocks)
			tt.prepareMocks(factory)
			handler := controller.NewApplyRetailerConfiguration(&controller.ApplyRetailerConfigurationDI{
				UC: uc,
			})

			applyOperation := handler.Operation()

			huma.Register(
				api,
				applyOperation,
				controller.ApplyRetailerConfiguration(handler),
			)

			bodyData, err := prepareApplyRetailerConfigurationRequest(tt.args.filePath)
			if err != nil {
				t.Fatal(err)
			}

			pathWithRetailerID := strings.ReplaceAll(applyOperation.Path, "{retailerId}", tt.args.pathParam)

			resp := api.Post(pathWithRetailerID, "X-Forwarded-Authorization: Bearer 1239hggasmmlnn8asdas1qsdasdas3441221", bodyData)
			assert.Equal(t, tt.expectedStatusCode, resp.Code)
		})
	}
}
