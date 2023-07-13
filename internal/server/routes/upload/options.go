package upload

import (
	"github.com/redhatinsights/insights-ingress-go/internal/storage"
	"github.com/redhatinsights/insights-ingress-go/internal/validators/kafka"
	"github.com/spf13/pflag"
)

type Options struct {
	DefaultMaxSize    int64             `mapstructure:"DefaultMaxSize"`
	MaxSizeMap        map[string]string `mapstructure:"MaxSizeMap"`
	MaxUploadMem      int64             `mapstructure:"MaxUploadMem"`
	BlackListedOrgIDs []string          `mapstructure:"BlackListedOrgIDs"`
	DebugUserAgent    string            `mapstructure:"DebugUserAgent"`

	Kafka   *kafka.Options   `mapstructure:"Kafka"`
	Storage *storage.Options `mapstructure:"Storage"`
}

// NewOptions returns an Options object with sensible defaults.
func NewOptions() *Options {
	return &Options{
		DefaultMaxSize:    100 * 1024 * 1024,
		MaxSizeMap:        map[string]string{},
		MaxUploadMem:      1024 * 1024 * 8,
		BlackListedOrgIDs: []string{},
		DebugUserAgent:    "",

		Kafka:   kafka.NewOptions(),
		Storage: storage.NewOptions(),
	}
}

// AddFlags sets up command line flags so a user can configure the application from the CLI.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}

	fs.Int64(prefix+"DefaultMaxSize", o.DefaultMaxSize, "the default max size of uploaded payloads")
	fs.StringToString(prefix+"MaxSizeMap", o.MaxSizeMap, "the max size map")
	fs.Int64(prefix+"MaxUploadMem", o.MaxUploadMem, "the default max size of uploaded payloads")
	fs.StringSlice(prefix+"BlackListedOrgIDs", o.BlackListedOrgIDs, "the black list of org ids")
	fs.String(prefix+"DebugUserAgent", o.DebugUserAgent, "the regex to match against user agents to debug")

	o.Kafka.AddFlags(fs, prefix+"Kafka")
	o.Storage.AddFlags(fs, prefix+"Storage")
}

// Complete adds values that haven't been supplied in some way to ensure the Options are complete.
func (o *Options) Complete() []error {
	var errs []error

	errs = append(errs, o.Kafka.Complete()...)
	errs = append(errs, o.Storage.Complete()...)

	return errs
}

// Validate checks that the values of the completed options are acceptable.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Kafka.Validate()...)
	errs = append(errs, o.Storage.Validate()...)

	return errs
}
