package kafka

import (
	"errors"

	"github.com/spf13/pflag"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

type Options struct {
	Brokers          []string `mapstructure:"Brokers"`
	GroupID          string   `mapstructure:"GroupID"`
	TrackerTopic     string   `mapstructure:"TrackerTopic"`
	DeliveryReports  bool     `mapstructure:"DeliveryReports"`
	AnnounceTopic    string   `mapstructure:"AnnounceTopic"`
	ValidUploadTypes []string `mapstructure:"ValidUploadTypes"`
	SecurityProtocol string   `mapstructure:"SecurityProtocol"`
	CA               string   `mapstructure:"CA"`
	Username         string   `mapstructure:"Username"`
	Password         string   `mapstructure:"Password"`
	SASLMechanism    string   `mapstructure:"SASLMechanism"`
}

func NewOptions() *Options {
	options := &Options{
		Brokers:          []string{"kafka:29092"},
		GroupID:          "ingress",
		TrackerTopic:     "platform.payload-status",
		DeliveryReports:  true,
		AnnounceTopic:    "platform.upload.announce",
		ValidUploadTypes: []string{"unit", "announce"},
		SecurityProtocol: "PLAINTEXT",
	}

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig
		broker := cfg.Kafka.Brokers[0]

		options.Brokers = clowder.KafkaServers
		options.TrackerTopic = clowder.KafkaTopics["platform.payload-status"].Name
		options.AnnounceTopic = clowder.KafkaTopics["platform.upload.announce"].Name

		if broker.SecurityProtocol != nil && *broker.SecurityProtocol != "" {
			options.SecurityProtocol = *broker.SecurityProtocol
		} else if broker.Sasl != nil && broker.Sasl.SecurityProtocol != nil && *broker.Sasl.SecurityProtocol != "" {
			options.SecurityProtocol = *broker.Sasl.SecurityProtocol
		}

		if broker.Authtype != nil {
			options.Username = *broker.Sasl.Username
			options.Password = *broker.Sasl.Password
			options.SASLMechanism = *broker.Sasl.SaslMechanism
		}

		if broker.Cacert != nil {
			caPath, err := cfg.KafkaCa(broker)
			if err != nil {
				panic("Kafka CA failed to write")
			}
			options.CA = caPath
		}
	}

	return options
}

// AddFlags sets up command line flags so a user can configure the application from the CLI.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefix string) {
	if prefix != "" {
		prefix = prefix + "."
	}
	fs.StringArray(prefix+"Brokers", o.Brokers, "the kafka brokers")
	fs.String(prefix+"GrouperID", o.GroupID, "the kafka group id")
	fs.String(prefix+"TrackerTopic", o.TrackerTopic, "the payload tracker topic")
	fs.Bool(prefix+"DeliveryReports", o.DeliveryReports, "delivery reports?")
	fs.String(prefix+"AnnounceTopic", o.AnnounceTopic, "the Kafka announce topic")
	fs.StringArray(prefix+"ValidUploadTypes", o.ValidUploadTypes, "the valid upload types")
	fs.String(prefix+"SecurityProtocol", o.SecurityProtocol, "the security protocol")
	fs.String(prefix+"CA", o.CA, "the CA")
	fs.String(prefix+"Username", o.Username, "the username")
	fs.String(prefix+"Password", o.Password, "the password")
	fs.String(prefix+"SASLMechanism", o.SASLMechanism, "the SASLMechanism")
}

// Complete adds values that haven't been supplied in some way to ensure the Options are complete.
func (o *Options) Complete() []error {
	var errs []error
	return errs
}

// Validate checks that the values of the completed options are acceptable.
func (o *Options) Validate() []error {
	var errs []error
	if len(o.Brokers) == 0 {
		errs = append(errs, errors.New("specify at least one kafka broker"))
	}
	return errs
}
