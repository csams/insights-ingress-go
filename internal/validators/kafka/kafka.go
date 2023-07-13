package kafka

import (
	"encoding/json"
	"errors"

	"github.com/redhatinsights/insights-ingress-go/internal/common"
	"github.com/redhatinsights/insights-ingress-go/internal/validators"
	"github.com/sirupsen/logrus"
)

// Validator posts requests to topics for validation
type Validator struct {
	CompletedConfig
}

// New constructs and initializes a new Kafka Validator
func New(c CompletedConfig, co common.CompletedConfig) (*Validator, error) {
	if producer, err := NewProducer(c, co); err != nil {
		return nil, err
	} else {
		go producer.Produce(c.ValidationProducerChannel, c.AnnounceTopic)
		return &Validator{
			c,
		}, nil
	}
}

// Validate validates a ValidationRequest
func (kv *Validator) Validate(vr *validators.Request) {
	data, err := json.Marshal(vr)
	if err != nil {
		kv.Log.WithFields(logrus.Fields{"error": err}).Error("failed to marshal json")
		return
	}
	kv.Log.WithFields(logrus.Fields{"data": data, "topic": kv.AnnounceTopic}).Debug("Posting data to topic")
	message := validators.ValidationMessage{
		Message: data,
		Headers: map[string]string{
			"service": vr.Service,
		},
	}
	if vr.Metadata.QueueKey != "" {
		message.Key = []byte(vr.Metadata.QueueKey)
	}

	kv.ValidationProducerChannel <- message
	incMessageProduced(vr.Service)
}

// ValidateService ensures that a service maps to a real topic
func (kv *Validator) ValidateService(service *validators.ServiceDescriptor) error {
	if kv.ValidUploadTypes.Contains(service.Service) {
		return nil
	}
	return errors.New("Upload type is not supported: " + service.Service)
}
