package logging

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"

	"github.com/redhatinsights/insights-ingress-go/internal/set"
)

// Common holds the options that are common to both options and Config.
type Common struct {
}

// Options contains the options to configure logging.
type Options struct {
	LogGroup           string `mapstructure:"logGroup"`
	AwsRegion          string `mapstructure:"AwsRegion"`
	AwsAccessKeyId     string `mapstructure:"AwsAccessKeyId"`
	AwsSecretAccessKey string `mapstructure:"AwsSecretAccessKey"`
	LogLevel string `mapstructure:"logLevel"`
}

// NewOptions returns an Options object with sensible defaults.
func NewOptions() *Options {
    options := &Options{
        LogGroup:           "platform-dev",
        AwsRegion:          "us-east-1",
        AwsAccessKeyId:     "",
        AwsSecretAccessKey: "",
		LogLevel: "INFO",
	}
	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig
		options.LogGroup = cfg.Logging.Cloudwatch.LogGroup
		options.AwsRegion = cfg.Logging.Cloudwatch.Region
		options.AwsAccessKeyId = cfg.Logging.Cloudwatch.AccessKeyId
		options.AwsSecretAccessKey = cfg.Logging.Cloudwatch.SecretAccessKey
	}
    return options
}

// AddFlags sets up command line flags so a user can configure the application from the CLI.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}
	fs.String(prefix+"LogGroup", o.LogGroup, "the log group")
	fs.String(prefix+"LogLevel", o.LogLevel, "the default log level")
	fs.String(prefix+"AwsRegion", o.AwsRegion, "the AWS region")
	fs.String(prefix+"AwsAccessKeyId", o.AwsAccessKeyId, "the AWS access key id")
	fs.String(prefix+"AwsSecretAccessKey", o.AwsSecretAccessKey, "the AWS secret access key")
}

// Complete adds values that haven't been supplied in some way to ensure the Options are complete.
func (o *Options) Complete() []error {
	if o.AwsAccessKeyId == "" {
		o.AwsAccessKeyId = os.Getenv("CW_AWS_ACCESS_KEY_ID")
	}

	if o.AwsSecretAccessKey == "" {
		o.AwsSecretAccessKey = os.Getenv("CW_AWS_SECRET_ACCESS_KEY")
	}
	return nil
}

// Validate checks that the values of the completed options are acceptable.
func (o *Options) Validate() []error {
	validLevels := set.New([]string{
		"",
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"FATAL",
    })
	var errs []error
	if !validLevels.Contains(o.LogLevel) {
		msg := fmt.Sprintf("invalid log level: %s", o.LogLevel)
		errs = append(errs, errors.New(msg))
	}
	return errs
}
