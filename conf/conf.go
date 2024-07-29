package conf

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

const (
	serviceName = "configurationManager"
)

// ServiceConfig represents the configuration of a service.
type ServiceConfig struct {
	Server      ServerConfig
	ServiceName string
	LMO         LMOService
}

type ServerConfig struct {
	BasePath  string `env:"BASE_PATH,default=/api/v3/configManager/"`
	Port      string `env:"PORT,default=8080"`
	ServerURL string `env:"SERVER_URL,default=http://localhost:8080"`
}

type LMOService struct {
	Host            string        `env:"BASE_URL,required"`
	Timeout         time.Duration `env:"TIMEOUT,default=15s"`
	ResponseTimeout time.Duration `env:"RESPONSE_TIMEOUT,default=30s"`
}

func New(ctx context.Context) (*ServiceConfig, error) {
	myConf := ServiceConfig{}
	if err := envconfig.Process(ctx, &myConf); err != nil {
		return nil, fmt.Errorf("failed to create service config: %w", err)
	}

	myConf.ServiceName = serviceName

	return &myConf, nil
}
