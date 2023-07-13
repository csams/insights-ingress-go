package storage

import (
	"github.com/minio/minio-go/v6"
)

type Config struct {
	*Options

	Client *minio.Client
}

type completedConfig struct {
	*Config
}

// CompletedConfig can be constructed only from Config.Complete.
type CompletedConfig struct {
	*completedConfig
}

// NewConfig creates a new configuration object from the given Options.
func NewConfig(o *Options) *Config {
	return &Config{
		Options: o,
	}
}

// Complete ensures the Config has all required values filled.
func (c *Config) Complete() (CompletedConfig, error) {
	if c.Client == nil {
		var err error
		if c.AwsRegion != "" {
			c.Client, err = minio.NewWithRegion(c.Endpoint, c.AwsAccessKeyId, c.AwsSecretAccessKey, c.UseSSL, c.AwsRegion)
		} else {
			c.Client, err = minio.New(c.Endpoint, c.AwsAccessKeyId, c.AwsSecretAccessKey, c.UseSSL)
		}

		if err != nil {
			return CompletedConfig{}, err
		}
	}
	return CompletedConfig{&completedConfig{
		c,
	}}, nil
}
