package server

import (
	"github.com/sirupsen/logrus"

	"github.com/redhatinsights/insights-ingress-go/internal/common"
	"github.com/redhatinsights/insights-ingress-go/internal/logging"
	"github.com/redhatinsights/insights-ingress-go/internal/server/routes/tracker"
	"github.com/redhatinsights/insights-ingress-go/internal/server/routes/upload"
	"github.com/redhatinsights/insights-ingress-go/internal/server/routes/version"
)

type Config struct {
	*Options

	Log *logrus.Logger

	Common  common.CompletedConfig
	Tracker *tracker.Config
	Upload  *upload.Config
	Version *version.Config
}

type completedConfig struct {
	*Config

	Tracker tracker.CompletedConfig
	Upload  upload.CompletedConfig
	Version version.CompletedConfig
}

type CompletedConfig struct {
	*completedConfig
}

// NewConfig creates a new configuration object from the given Options.
func NewConfig(o *Options, commonConfig common.CompletedConfig) *Config {
	return &Config{
		Options: o,

		Common:  commonConfig,
		Tracker: tracker.NewConfig(o.Tracker, commonConfig),
		Upload:  upload.NewConfig(o.Upload, commonConfig),
		Version: version.NewConfig(o.Version),
	}
}

func (c *Config) Complete() (CompletedConfig, error) {
	trackerCompleted, err := c.Tracker.Complete()
	if err != nil {
		return CompletedConfig{}, err
	}

	uploadCompleted, err := c.Upload.Complete()
	if err != nil {
		return CompletedConfig{}, err
	}

	versionCompleted := c.Version.Complete()

	if c.Log == nil {
		c.Log = logging.Log
	}

	return CompletedConfig{&completedConfig{
		Config: c,

		Tracker: trackerCompleted,
		Upload:  uploadCompleted,
		Version: versionCompleted,
	}}, nil
}
