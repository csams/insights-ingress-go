package kafka

import (
	"time"

	"github.com/redhatinsights/insights-ingress-go/internal/common"
	"github.com/redhatinsights/insights-ingress-go/internal/errors"
	l "github.com/redhatinsights/insights-ingress-go/internal/logging"
	"github.com/redhatinsights/insights-ingress-go/internal/validators"
	"github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	prom "github.com/prometheus/client_golang/prometheus"
)

type Producer struct {
	CompletedConfig

	Log      *logrus.Logger
	Producer *kafka.Producer
}

func newConfigMap(config CompletedConfig, commonConfig common.CompletedConfig) (kafka.ConfigMap, error) {
	configMap := kafka.ConfigMap{
		"bootstrap.servers":   config.Brokers[0],
		"go.delivery.reports": config.DeliveryReports,
	}

	var errs []error
	if commonConfig.Debug {
		if err := configMap.SetKey("debug", "protocol,broker,topic"); err != nil {
			errs = append(errs, err)
		}
	}

	if config.CA != "" {
		if err := configMap.SetKey("ssl.ca.location", config.CA); err != nil {
			errs = append(errs, err)
		}
		if err := configMap.SetKey("security.protocol", config.SecurityProtocol); err != nil {
			errs = append(errs, err)
		}

		if config.SASLMechanism != "" {
			if err := configMap.SetKey("sasl.mechanism", config.SASLMechanism); err != nil {
				errs = append(errs, err)
			}
			if err := configMap.SetKey("sasl.username", config.Username); err != nil {
				errs = append(errs, err)
			}
			if err := configMap.SetKey("sasl.password", config.Password); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return kafka.ConfigMap{}, errors.NewAggregate(errs)
	}

	return configMap, nil
}

func NewProducer(config CompletedConfig, commonConfig common.CompletedConfig) (*Producer, error) {
	if config.Log == nil {
		config.Log = l.Log
	}

	var configMap kafka.ConfigMap
	if c, err := newConfigMap(config, commonConfig); err == nil {
		configMap = c
	} else {
		config.Log.WithFields(logrus.Fields{"error": err}).Error("Error creating kafka producer")
		return nil, err
	}

	if p, err := kafka.NewProducer(&configMap); err == nil {
		return &Producer{
			CompletedConfig: config,
			Producer:        p,
		}, nil

	} else {
		config.Log.WithFields(logrus.Fields{"error": err}).Error("Error creating kafka producer")
		return nil, err
	}
}

// Produce consumes in and produces to the topic in config
// Each message is sent to the writer via a goroutine so that the internal batch
// buffer has an opportunity to fill.
func (config *Producer) Produce(in chan validators.ValidationMessage, topic string) {

	defer config.Producer.Close()

	for v := range in {
		go func(v validators.ValidationMessage) {
			delivery_chan := make(chan kafka.Event)
			defer close(delivery_chan)
			producerCount.Inc()
			defer producerCount.Dec()
			start := time.Now()
			kafkaHeaders := make([]kafka.Header, len(v.Headers))
			i := 0
			for key, value := range v.Headers {
				kafkaHeaders[i] = kafka.Header{
					Key:   key,
					Value: []byte(value),
				}
				i++
			}
			config.Producer.Produce(&kafka.Message{
				Headers: kafkaHeaders,
				TopicPartition: kafka.TopicPartition{
					Topic:     &topic,
					Partition: kafka.PartitionAny,
				},
				Value: v.Message,
				Key:   v.Key,
			}, delivery_chan)
			messagePublishElapsed.With(prom.Labels{"topic": topic}).Observe(time.Since(start).Seconds())

			e := <-delivery_chan
			m := e.(*kafka.Message)

			if m.TopicPartition.Error != nil {
				l.Log.WithFields(logrus.Fields{"error": m.TopicPartition.Error}).Error("Error publishing to kafka")
				in <- v
				publishFailures.With(prom.Labels{"topic": topic}).Inc()
				return
			} else {
				messagesPublished.With(prom.Labels{"topic": topic}).Inc()
			}
		}(v)
	}
}
