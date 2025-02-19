package kafka

import (
	"context"
	"fmt"
	"log"

	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

type KafkaWriterInterface interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type MockWriter struct {
	mock.Mock
}

func (w *MockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := w.Called(ctx, msgs)
	return args.Error(0)
}

func NewConsumer(cfg *viper.Viper, topic string) *kafka.Reader {
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.GetStringSlice("Kafka_Brokers"),
		Topic:       topic,
		GroupID:     cfg.GetString("Kafka_Group_ID"),
		StartOffset: cfg.GetInt64("Kafka_Consumer_Offset"),
	})

	return consumer
}

func NewProducer(cfg *viper.Viper, topic string) *kafka.Writer {
	producer := &kafka.Writer{
		Addr:  kafka.TCP(cfg.GetStringSlice("Kafka_Brokers")[0]),
		Topic: topic,
	}

	return producer
}

func NewConsumerEventLoop(
	ctx context.Context,
	consumer *kafka.Reader,
	handler func(context.Context, kafka.Message),
	errors chan<- error,
) (start func()) {
	return func() {
		for {
			m, err := consumer.ReadMessage(ctx)
			if err != nil {
				log.Println(err)
				errors <- err
			}
			handler(ctx, m)
		}
	}
}

func GetHeader(msg kafka.Message, key string) (string, error) {
	for _, value := range msg.Headers {
		if value.Key == key {
			return string(value.Value), nil
		}
	}

	return "", fmt.Errorf("Header not found: %s", key)
}
