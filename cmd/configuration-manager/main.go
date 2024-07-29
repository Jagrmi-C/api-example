package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"google.golang.org/api/idtoken"

	"gitlab.com/jc88/api-example/conf"
	"gitlab.com/jc88/api-example/infra/adapter/configurationlmo"
	"gitlab.com/jc88/api-example/infra/controller"
	"gitlab.com/jc88/api-example/infra/errorx"
	"gitlab.com/jc88/api-example/infra/gohttp"
	"gitlab.com/jc88/api-example/usecase"
)

const (
	// pathToOpenAPIFile is the path to openapi file.
	pathToOpenAPIFile = "../../docs/openapi.yaml"
	// filePerms is the value needed to write bytes in file.
	filePerms = 0600 //nolint:all
)

func handleErr(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	configSvc, err := conf.New(ctx)
	if err != nil {
		handleErr(err)
	}

	svcLog, err := setUpLogger(configSvc)
	if err != nil {
		handleErr(err)
	}

	router := mux.NewRouter()
	config := createNewOpenAPIConfig(configSvc)

	api := humamux.New(router, config)

	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		setCustomError()
		addRoutes(ctx, svcLog, api, configSvc)

		hooks.OnStart(func() {
			_, _ = fmt.Printf("Starting server on port %d...\n", options.Port)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router); err != nil {
				handleErr(err)
			}
		})
	})

	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			// Use downgrade to return OpenAPI 3.0.3 YAML since oapi-codegen doesn't
			// support OpenAPI 3.1 fully yet. Use `.YAML()` instead for 3.1.
			openAPIData, err := api.OpenAPI().DowngradeYAML()
			if err != nil {
				handleErr(err)
			}

			if err := os.WriteFile(pathToOpenAPIFile, openAPIData, filePerms); err != nil {
				handleErr(err)
			}
		},
	})

	cli.Run()
}

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8080"`
}

func setUpLogger(cfg *conf.ServiceConfig) (*slog.Logger, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Logger Initialized", "count", 3)

	return logger, nil
}

func setCustomError() {
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		entities := make([]errorx.ErrorEntity, 0, len(errs))

		for _, err := range errs {
			errID, _ := uuid.NewV4()

			errEnt := errorx.ErrorEntity{
				ID:     errID.String(),
				Status: status,
			}

			switch errT := err.(type) {
			case *huma.ErrorDetail:
				result := strings.Split(errT.Location, ".")
				errEnt.Detail = errT.Message

				if len(result) > 0 {
					switch result[0] {
					case "body":
						errEnt.Source = &errorx.Source{
							Pointer: errT.Location,
						}
					case "path":
						errEnt.Source = &errorx.Source{
							Path: errT.Location,
						}
					case "query":
						errEnt.Source = &errorx.Source{
							Parameter: errT.Location,
						}
					case "header":
						errEnt.Source = &errorx.Source{
							Header: errT.Location,
						}
					default:
						errEnt.Source = &errorx.Source{
							Pointer: errT.Location,
						}
					}
				}
			case *errorx.Error:
				errEnt.Detail = errT.Details()
				if errEnt.Detail == "" {
					errEnt.Detail = message
				}
				errEnt.Status = status
				errEnt.Title = errT.Code().ErrMsg()
			default:
				errEnt.Detail = message
			}

			entities = append(entities, errEnt)
		}

		return &errorx.ErrorsRespo{
			Status: status,
			Errors: entities,
		}
	}
}

func addRoutes(ctx context.Context, svcLog *slog.Logger, api huma.API, configSvc *conf.ServiceConfig) {
	preparedJSONHeaders := make(http.Header, 1)
	preparedJSONHeaders.Add(gohttp.HeaderContentType, gohttp.ContentTypeJSON)

	client, err := idtoken.NewClient(ctx, configSvc.LMO.Host)
	if err != nil {
		handleErr(err)
	}

	defaultJsonClient := gohttp.NewClientBuilder().
		SetHeaders(preparedJSONHeaders).
		SetConnectionTimeOut(configSvc.LMO.Timeout).
		SetResponseHeaderTimeout(configSvc.LMO.ResponseTimeout).
		SetHTTPClient(client).
		Build()

	lmoAdapter := configurationlmo.NewRetailerConfigurationLMOAdapter(
		configSvc.LMO.Host,
		defaultJsonClient,
	)

	svcLog.Info(fmt.Sprintf("Adapter REST HTTP to LMO was initialized: %s", configSvc.LMO.Host))

	uc := usecase.NewRetailerConfigurationsApplier(
		&usecase.RetailerConfigurationsApplierDI{
			LMOAdapter: lmoAdapter,
			Logger:     svcLog,
		},
	)

	applyRetailerConfigsHTTPX := controller.NewApplyRetailerConfiguration(&controller.ApplyRetailerConfigurationDI{
		UC: uc,
	})

	huma.Register(
		api,
		applyRetailerConfigsHTTPX.Operation(),
		controller.ApplyRetailerConfiguration(applyRetailerConfigsHTTPX),
	)
}

func createNewOpenAPIConfig(configSvc *conf.ServiceConfig) huma.Config {
	config := huma.DefaultConfig("Configurations manager API", "1.0.0")

	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	config.Tags = []*huma.Tag{
		{
			Name:        "Configurations",
			Description: "Managing retailer configurations",
		},
	}

	config.Info.Contact = &huma.Contact{
		Name:  "test",
		URL:   "https://test.co/contact-us/",
		Email: "test@test.com",
	}

	config.Info.Description = "API for synchronization REMS configurations and LMO"
	config.Info.License = &huma.License{
		Name: "internal",
		URL:  "https://test.co",
	}

	config.Servers = append(
		config.Servers,
		&huma.Server{
			URL:         configSvc.Server.ServerURL,
			Description: "Current configuration server URL",
		},
		&huma.Server{
			URL:         "http://localhost:8080",
			Description: "For local development",
		},
	)

	return config
}
