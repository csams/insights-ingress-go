package version

import (
	"github.com/redhatinsights/insights-ingress-go/internal/logging"
	"github.com/sirupsen/logrus"
)

type Config struct {
	*Options

	Log *logrus.Logger

	IngressVersion *IngressVersion
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	*completedConfig
}

// NewConfig creates a new configuration object from the given Options.
func NewConfig(o *Options) *Config {
	return &Config{
		Options: o,
	}
}

func (c *Config) Complete() CompletedConfig {
	if c.IngressVersion == nil {
		c.IngressVersion = &IngressVersion{
			Commit:  c.OpenshiftBuildCommit,
			Version: c.Version,
		}
	}

	if c.Log == nil {
		c.Log = logging.Log
	}

	return CompletedConfig{&completedConfig{
		Config: c,
	}}
}
