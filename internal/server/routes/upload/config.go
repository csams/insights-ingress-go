package upload

import (
	"github.com/redhatinsights/insights-ingress-go/internal/common"
	"github.com/redhatinsights/insights-ingress-go/internal/logging"
	"github.com/redhatinsights/insights-ingress-go/internal/storage"
	"github.com/redhatinsights/insights-ingress-go/internal/validators/kafka"
	"github.com/sirupsen/logrus"
)

type Config struct {
	*Options

	Log     *logrus.Logger

	Common  common.CompletedConfig
	Kafka   *kafka.Config
	Storage *storage.Config
}

type completedConfig struct {
	*Config

	Kafka   kafka.CompletedConfig
	Storage storage.CompletedConfig
}

type CompletedConfig struct {
	*completedConfig
}

// NewConfig creates a new configuration object from the given Options.
func NewConfig(o *Options, commonConfig common.CompletedConfig) *Config {
	return &Config{
		Options: o,

		Common:  commonConfig,
		Kafka:   kafka.NewConfig(o.Kafka),
		Storage: storage.NewConfig(o.Storage),
	}
}

func (c *Config) Complete() (CompletedConfig, error) {
	if c.Log == nil {
		c.Log = logging.Log
	}

	storageCompleted, err := c.Storage.Complete()
	if err != nil {
		return CompletedConfig{}, err
	}

	return CompletedConfig{&completedConfig{
		Config:  c,
		Kafka:   c.Kafka.Complete(),
		Storage: storageCompleted,
	}}, nil
}
