package logging

import (
	"flag"

	"github.com/redhatinsights/insights-ingress-go/internal/common"
	"github.com/sirupsen/logrus"
)

type Config struct {
    *Options
    LogLevel logrus.Level

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
    var logLevel logrus.Level
    switch o.LogLevel {
    case "DEBUG":
        logLevel = logrus.DebugLevel
    case "ERROR":
        logLevel = logrus.ErrorLevel
    default:
        logLevel = logrus.InfoLevel
    }
    if flag.Lookup("test.v") != nil {
        logLevel = logrus.FatalLevel
    }

	return &Config{
        Options: o,
        Common: commonConfig,
        LogLevel: logLevel,
	}
}

// Complete ensures the Config has all required values filled.
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{&completedConfig{
        c,
	}}
}
