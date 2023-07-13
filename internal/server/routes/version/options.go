package version

import (
	"os"

	"github.com/spf13/pflag"
)

type Options struct {
	OpenshiftBuildCommit string `mapstructure:"OpenshiftBuildCommit"`
	Version              string `mapstructure:"Version"`
}

// NewOptions returns an Options object with sensible defaults.
func NewOptions() *Options {
	options := &Options{
		OpenshiftBuildCommit: "notrunninginopenshift",
		Version:              os.Getenv("1.0.8"),
	}
	return options
}

// AddFlags sets up command line flags so a user can configure the application from the CLI.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}

	fs.String(prefix+"OpenshiftBuildCommit", o.OpenshiftBuildCommit, "the openshift build commit")
	fs.String(prefix+"Version", o.Version, "the ingress server version")
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
