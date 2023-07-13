package kafka

import (
	"time"

	p "github.com/prometheus/client_golang/prometheus"
	pa "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	payloadsProcessed = pa.NewCounterVec(p.CounterOpts{
		Name: "ingress_processed_payloads",
		Help: "The total number of processed events",
	}, []string{"outcome"})

	validationElapsed = pa.NewHistogramVec(p.HistogramOpts{
		Name: "ingress_validate_elapsed_seconds",
		Help: "Number of seconds spent to validating",
	}, []string{"outcome"})

	messageProduced = pa.NewCounterVec(p.CounterOpts{
		Name: "ingress_message_produced",
		Help: "The total number of messages produced",
	}, []string{"service"})

	messagesPublished = pa.NewCounterVec(p.CounterOpts{
		Name: "ingress_kafka_produced",
		Help: "Number of messages produced to kafka",
	}, []string{"topic"})
	messagePublishElapsed = pa.NewHistogramVec(p.HistogramOpts{
		Name: "ingress_publish_seconds",
		Help: "Number of seconds spent writing kafka messages",
	}, []string{"topic"})
	publishFailures = pa.NewCounterVec(p.CounterOpts{
		Name: "ingress_kafka_produce_failures",
		Help: "Number of times a message was failed to be produced",
	}, []string{"topic"})
	producerCount = pa.NewGauge(p.GaugeOpts{
		Name: "ingress_kafka_producer_go_routine_count",
		Help: "Number of go routines currently publishing to kafka",
	})
)

func inc(outcome string) {
	payloadsProcessed.With(p.Labels{"outcome": outcome}).Inc()
}

func incMessageProduced(service string) {
	messageProduced.With(p.Labels{"service": service}).Inc()
}

func observeValidationElapsed(timestamp time.Time, outcome string) {
	validationElapsed.With(p.Labels{
		"outcome": outcome,
	}).Observe(time.Since(timestamp).Seconds())
}
