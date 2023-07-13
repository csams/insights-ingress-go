package tracker

import (
	"time"

	"github.com/spf13/pflag"
)

type Options struct {
	Url               string
	HTTPClientTimeout time.Duration
	Enabled           bool
}

// NewOptions returns an Options object with sensible defaults.
func NewOptions() *Options {
	return &Options{
		Url:               "http://payload-tracker/v1/payloads/",
		Enabled:           true,
		HTTPClientTimeout: 10 * time.Second,
	}
}

// AddFlags sets up command line flags so a user can configure the application from the CLI.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}

	fs.String(prefix+"Url", o.Url, "the payload tracker URL")
	fs.Bool(prefix+"Enabled", o.Enabled, "is auth enabled")
	fs.Duration(prefix+"HTTPClientTimeout", o.HTTPClientTimeout, "http client timeout")
}

// Complete adds values that haven't been supplied in some way to ensure the Options are complete.
func (o *Options) Complete() []error {
	var errs []error
	return errs
}

// Validate checks that the values of the completed options are acceptable.
func (o *Options) Validate() []error {
	var errs []error
	return errs
}
