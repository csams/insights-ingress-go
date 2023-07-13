package kafka

import (
	"github.com/redhatinsights/insights-ingress-go/internal/logging"
	"github.com/redhatinsights/insights-ingress-go/internal/set"
	"github.com/redhatinsights/insights-ingress-go/internal/validators"
	"github.com/sirupsen/logrus"
)

// Config configures a new Kafka Validator
type Config struct {
    *Options

    Log *logrus.Logger

    ValidationProducerChannel chan validators.ValidationMessage
}

type completedConfig struct {
	*Config

    ValidUploadTypes set.Set[string]
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
    if c.Log == nil {
        c.Log = logging.Log
    }

    if c.ValidationProducerChannel == nil {
        c.ValidationProducerChannel = make(chan validators.ValidationMessage, 100)
    }

	return CompletedConfig{&completedConfig{
        Config: c,
        ValidUploadTypes: set.New(c.ValidUploadTypes),
	}}
}
