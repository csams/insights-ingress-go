package common

import (
	"os"

	"github.com/spf13/pflag"
)

type Options struct {
	Auth     bool   `mapstructure:"Auth"`
	Debug    bool   `mapstructure:"Debug"`
	Hostname string `mapstructure:"Hostname"`
}

// NewOptions returns an Options object with sensible defaults.
func NewOptions() *Options {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return &Options{
		Auth:     true,
		Debug:    false,
		Hostname: hostname,
	}
}

// AddFlags sets up command line flags so a user can configure the application from the CLI.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}
	fs.Bool(prefix+"Auth", o.Auth, "enable authentation")
	fs.Bool(prefix+"Debug", o.Debug, "debug mode?")
	fs.String(prefix+"Hostname", o.Hostname, "the hostname")
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
