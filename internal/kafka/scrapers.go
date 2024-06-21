package kafka

import (
	"bytes"
	"context"
	"encoding/json"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/getsentry/sentry-go"
)

type ScraperMessageBroker struct {
	producer *kafka.Producer
}

type event struct {
	Type  string
	Value internal.ScrapeRequest
}

func NewScraperMessageBroker(producer *kafka.Producer) *ScraperMessageBroker {
	return &ScraperMessageBroker{
		producer: producer,
	}
}

func (s *ScraperMessageBroker) Created(ctx context.Context, params internal.ScrapeRequest) error {
	return s.publish(ctx, "ScraperMessageBroker.Create", "scrape-request", string(params.Type), params)
}

func (s *ScraperMessageBroker) publish(ctx context.Context, spanName, topic, message string, params internal.ScrapeRequest) error {
	defer newSentrySpan(ctx, spanName).Finish()

	var b bytes.Buffer

	evt := event{
		Type:  message,
		Value: params,
	}

	if err := json.NewEncoder(&b).Encode(evt); err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "json.Encode")
	}

	if err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(params.ID),
		Value: b.Bytes(),
	}, nil); err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "product.Producer")
	}

	return nil
}

func newSentrySpan(ctx context.Context, operation string) *sentry.Span {
	span := sentry.StartSpan(ctx, operation)
	span.Name = "fourleaves.studio/manga-scraper/internal/kafka"

	span.SetTag("messaging.system", "kafka")
	span.SetTag("messaging.destination", "scrape-request")

	return span
}
