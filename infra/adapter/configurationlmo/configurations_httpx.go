package configurationlmo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"gitlab.com/jc88/api-example/domain"
	"gitlab.com/jc88/api-example/infra/errorx"
	"gitlab.com/jc88/api-example/infra/gohttp"
)

const (
	retrieveInstructionsMsg = "retrieve instructions by retailer"
	unauthorizedAccessMsg   = "unauthorized access"
	unauthorizedAccessToMsg = "unauthorized access to %s"
	connectionErrMsg        = "unable to connect to API - %s: %w"
	createConfigurationMsg  = "process to creation retailer configuration in LMO"
	removeConfigurationMsg  = "process to remove retailer configuration in LMO"
	unmarshalRespoBody      = "unmarshal response body"

	// clientResponseWithDetailsMsg is used as template for error messages.
	clientResponseWithDetailsMsg = "%s: status code: %d"
)

// HeaderAuthorizationForwarded is a header that needed to call LMO.
const HeaderAuthorizationForwarded = "X-Forwarded-Authorization"

var (
	ErrInvalidAuthorizationRequest = errors.New("invalid auth0 authorization request")
	ErrAuth0ConfigurationRequired  = errors.New("auth0 configuration is required")
	ErrValidationLMORequest        = errors.New("request to LMO system returned invalid client status code")
	ErrConnectionLMOEndpoint       = errors.New("connection to LMO system returned internal error server status code")
)

// HTTPResponseDebug represents data that should be added to the log payload in Debug field.
type HTTPResponseDebug struct {
	ResponseBody       any
	ResourceURL        string
	ResponseStatusCode int
}

type configurationManagerLMOAdapter struct {
	httpClient gohttp.Client
	baseURL    string
}

func NewRetailerConfigurationLMOAdapter(baseURL string, client gohttp.Client) *configurationManagerLMOAdapter {
	adapt := configurationManagerLMOAdapter{
		baseURL:    baseURL,
		httpClient: client,
	}

	return &adapt
}

func (u *configurationManagerLMOAdapter) PrepareEndpoint(method string, retailerID uuid.UUID) string {
	switch method {
	case "GET", "POST":
		return fmt.Sprintf("%s/retailers/%s/configurations", u.baseURL, retailerID.String())
	case "DELETE":
		return fmt.Sprintf("%s/retailers/%s", u.baseURL, retailerID.String())
	default:
		return ""
	}
}

func (u *configurationManagerLMOAdapter) Create(ctx context.Context, config *domain.RetailerConfiguration, auth0token string) error {
	reqStruct, err := NewRetailerConfigurationCreateReq(config)
	if err != nil {
		return fmt.Errorf("unable to prepare request to API: %w", err)
	}

	extraHeader := http.Header{}
	extraHeader.Add(HeaderAuthorizationForwarded, auth0token)

	req := gohttp.RequestAttributes{
		Method:  http.MethodGet,
		URL:     u.PrepareEndpoint(http.MethodPost, config.RetailerID),
		Headers: extraHeader,
		Body:    reqStruct,
		Retry:   false,
	}

	resp, err := u.httpClient.DoReq(ctx, &req)
	if err != nil {
		return fmt.Errorf("unable to connect to API: %s - %w", u.baseURL, err)
	}

	switch resp.StatusCode() {
	case http.StatusCreated:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return errorx.NewErrorf(
			errorx.ErrCodeAuth0AuthMalformed,
			fmt.Sprintf(unauthorizedAccessToMsg, req.URL),
		)
	case http.StatusConflict:
		return errorx.NewErrorf(
			errorx.ErrCodeRetailerConfigurationAlreadyExist,
			fmt.Sprintf(clientResponseWithDetailsMsg, createConfigurationMsg, resp.StatusCode()),
		)
	case http.StatusBadRequest, http.StatusNotFound:
		return errorx.WrapErrorfWithDebug(
			ErrValidationLMORequest,
			resp.String(),
			errorx.ErrorCodeInputValidation,
			fmt.Sprintf(clientResponseWithDetailsMsg, createConfigurationMsg, resp.StatusCode()),
		)
	default:
		return errorx.WrapErrorfWithDebug(
			ErrConnectionLMOEndpoint,
			resp.String(),
			errorx.ErrCodeInternalError,
			fmt.Sprintf(clientResponseWithDetailsMsg, createConfigurationMsg, resp.StatusCode()),
		)
	}
}

func (u *configurationManagerLMOAdapter) GetConfigurationsByRetailerID(ctx context.Context, retailerID uuid.UUID) (*domain.RetailerConfiguration, error) {
	extraHeader := http.Header{}

	req := gohttp.RequestAttributes{
		Method:  http.MethodGet,
		URL:     u.PrepareEndpoint(http.MethodGet, retailerID),
		Headers: extraHeader,
		Retry:   false,
	}

	resp, err := u.httpClient.DoReq(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf(connectionErrMsg, u.baseURL, err)
	}

	respo := GetConfigurationsByRetailerRespo{}
	if err := resp.UnMarshalJson(&respo); err != nil {
		return nil, errorx.WrapErrorfWithDebug(
			err,
			HTTPResponseDebug{
				ResponseStatusCode: resp.StatusCode(),
				ResponseBody:       resp.String(),
			},
			errorx.ErrCodeInternalError,
			unmarshalRespoBody,
		)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		retailerId, err := uuid.FromString(respo.Data.RetailerID)
		if err != nil {
			return nil, errorx.WrapErrorfWithDebug(
				err,
				HTTPResponseDebug{
					ResponseStatusCode: resp.StatusCode(),
					ResponseBody:       resp.String(),
				},
				errorx.ErrCodeInternalError,
				retrieveInstructionsMsg,
			)
		}

		instructions := domain.RetailerConfiguration{
			RetailerID: retailerId,
		}
		return &instructions, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, errorx.NewErrorf(
			errorx.ErrCodeAuth0AuthMalformed,
			fmt.Sprintf(unauthorizedAccessToMsg, req.URL),
		)
	case http.StatusBadRequest, http.StatusNotFound, http.StatusConflict:
		return nil, errorx.WrapErrorfWithDebug(
			ErrValidationLMORequest,
			resp,
			errorx.ErrorCodeInputValidation,
			fmt.Sprintf(clientResponseWithDetailsMsg, retrieveInstructionsMsg, resp.StatusCode()),
		)
	default:
		return nil, errorx.WrapErrorfWithDebug(
			ErrConnectionLMOEndpoint,
			resp,
			errorx.ErrCodeInternalError,
			fmt.Sprintf(clientResponseWithDetailsMsg, retrieveInstructionsMsg, resp.StatusCode()),
		)
	}
}

func (u *configurationManagerLMOAdapter) DeleteConfigurationsByRetailerID(ctx context.Context, retailerID uuid.UUID, auth0token string) error {
	extraHeader := http.Header{}
	extraHeader.Add(HeaderAuthorizationForwarded, auth0token)

	req := gohttp.RequestAttributes{
		Method:  http.MethodDelete,
		URL:     u.PrepareEndpoint(http.MethodDelete, retailerID),
		Headers: extraHeader,
		Retry:   false,
	}

	resp, err := u.httpClient.DoReq(ctx, &req)
	if err != nil {
		return fmt.Errorf(connectionErrMsg, u.baseURL, err)
	}

	switch resp.StatusCode() {
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return errorx.NewErrorf(
			errorx.ErrCodeAuth0AuthMalformed,
			fmt.Sprintf(unauthorizedAccessToMsg, req.URL),
		)
	case http.StatusBadRequest, http.StatusNotFound, http.StatusConflict:
		return errorx.WrapErrorfWithDebug(
			ErrValidationLMORequest,
			resp,
			errorx.ErrorCodeInputValidation,
			fmt.Sprintf(clientResponseWithDetailsMsg, removeConfigurationMsg, resp.StatusCode()),
		)
	default:
		return errorx.WrapErrorfWithDebug(
			ErrConnectionLMOEndpoint,
			resp,
			errorx.ErrCodeInternalError,
			fmt.Sprintf(clientResponseWithDetailsMsg, removeConfigurationMsg, resp.StatusCode()),
		)
	}
}
