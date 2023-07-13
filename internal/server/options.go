package server

import (
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"

	"github.com/redhatinsights/insights-ingress-go/internal/server/routes/tracker"
	"github.com/redhatinsights/insights-ingress-go/internal/server/routes/upload"
	"github.com/redhatinsights/insights-ingress-go/internal/server/routes/version"
	"github.com/spf13/pflag"
)

type Options struct {
	PublicPort  int    `mapstructure:"PublicPort"`
	MetricsPort int    `mapstructure:"MetricsPort"`
	Profile     bool   `mapstructure:"Profile"`
	TlsCAPath   string `mapstructure:"TlsCAPath"`

	Tracker *tracker.Options `mapstructure:"Track"`
	Upload  *upload.Options  `mapstructure:"Options"`
	Version *version.Options `mapstructure:"Version"`
}

// NewOptions returns an Options object with sensible defaults.
func NewOptions() *Options {
	options := &Options{
		PublicPort:  3000,
		MetricsPort: 8080,
		Profile:     false,
		TlsCAPath:   "",

		Tracker: tracker.NewOptions(),
		Upload:  upload.NewOptions(),
		Version: version.NewOptions(),
	}

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig

		options.TlsCAPath = *cfg.TlsCAPath
		options.PublicPort = *cfg.PublicPort
		options.MetricsPort = cfg.MetricsPort
	}

	return options
}

// AddFlags sets up command line flags so a user can configure the application from the CLI.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}

	fs.Int(prefix+"PublicPort", o.PublicPort, "the public port ingress will listen on")
	fs.Int(prefix+"MetricsPort", o.PublicPort, "the port the metrics server will listen on")
	fs.Bool(prefix+"Profile", o.Profile, "start the profiler?")
	fs.String(prefix+"TlsCAPath", o.TlsCAPath, "The TLS CA path")

	o.Tracker.AddFlags(fs, prefix+"Tracker")
	o.Upload.AddFlags(fs, prefix+"Upload")
	o.Version.AddFlags(fs, prefix+"Version")
}

// Complete adds values that haven't been supplied in some way to ensure the Options are complete.
func (o *Options) Complete() []error {
	var errs []error

	errs = append(errs, o.Tracker.Complete()...)
	errs = append(errs, o.Upload.Complete()...)
	errs = append(errs, o.Version.Complete()...)

	return errs
}

// Validate checks that the values of the completed options are acceptable.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Tracker.Validate()...)
	errs = append(errs, o.Upload.Validate()...)
	errs = append(errs, o.Version.Validate()...)

	return errs
}
