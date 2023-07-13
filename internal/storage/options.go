package storage

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

type Options struct {
	StageBucket        string `mapstructure:"StageBucket"`
	UseSSL             bool   `mapstructure:"UseSSL"`
	Endpoint           string `mapstructure:"Endpoint"`
	AwsAccessKeyId     string `mapstructure:"AwsAccessKeyId"`
	AwsSecretAccessKey string `mapstructure:"AwsSecretAccessKey"`
	AwsRegion          string `mapstructure:"AwsRegion"`
}

// NewOptions returns an Options object with sensible defaults.
func NewOptions() *Options {
	options := &Options{
		StageBucket: "available",
		Endpoint:    "s3.amazonaws.com",
	}

	if clowder.IsClowderEnabled() {
		sb := os.Getenv("INGRESS_STAGEBUCKET")
		if bucket, found := clowder.ObjectBuckets[sb]; found {
			options.StageBucket = bucket.Name
		}

		cfg := clowder.LoadedConfig

		options.Endpoint = fmt.Sprintf("%s:%d", cfg.ObjectStore.Hostname, cfg.ObjectStore.Port)
		options.UseSSL = cfg.ObjectStore.Tls

		if cfg.ObjectStore.Buckets[0].AccessKey != nil {
			options.AwsAccessKeyId = *cfg.ObjectStore.Buckets[0].AccessKey
		}
		if cfg.ObjectStore.Buckets[0].SecretKey != nil {
			options.AwsSecretAccessKey = *cfg.ObjectStore.Buckets[0].SecretKey
		}
		if cfg.ObjectStore.Buckets[0].Region != nil {
			options.AwsRegion = *cfg.ObjectStore.Buckets[0].Region
		}
	}

	return options
}

// AddFlags sets up command line flags so a user can configure the application from the CLI.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}
	fs.String(prefix+"StageBucket", o.StageBucket, "the stage bucket")
	fs.Bool(prefix+"UseSSL", o.UseSSL, "use SSL")
	fs.String(prefix+"Endpoint", o.StageBucket, "the object storage endpoint")
	fs.String(prefix+"AwsRegion", o.AwsRegion, "the AWS region")
	fs.String(prefix+"AwsAccessKeyId", o.AwsAccessKeyId, "the AWS access key id")
	fs.String(prefix+"AwsSecretAccessKey", o.AwsSecretAccessKey, "the AWS secret access key")
}

// Complete adds values that haven't been supplied in some way to ensure the Options are complete.
func (o *Options) Complete() []error {
	return nil
}

// Validate checks that the values of the completed options are acceptable.
func (o *Options) Validate() []error {
	var errs []error
	if o.StageBucket == "" {
		errs = append(errs, errors.New("StageBucket is empty"))
	}
	return errs
}
