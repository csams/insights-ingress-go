package tracker

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/redhatinsights/insights-ingress-go/internal/common"
	"github.com/redhatinsights/insights-ingress-go/internal/logging"
)

type Config struct {
	*Options

	Client *http.Client
	Log    *logrus.Logger

	Common common.CompletedConfig
}

type completedConfig struct {
	*Config
}

// CompletedConfig can be constructed only from Config.Complete.
type CompletedConfig struct {
	*completedConfig
}

// NewConfig creates a new configuration object from the given Options.
func NewConfig(o *Options, commonConfig common.CompletedConfig) *Config {
	return &Config{
		Options: o,
        Common: commonConfig,
	}
}

// Complete ensures the Config has all required values filled.
func (c *Config) Complete() (CompletedConfig, error) {
	if c.Client == nil {
		c.Client = &http.Client{
			Timeout: c.HTTPClientTimeout,
		}
	}

	if c.Log == nil {
		c.Log = logging.Log
	}

	return CompletedConfig{&completedConfig{
		c,
	}}, nil
}
