package kafkaservice

import (
	"github.com/segmentio/kafka-go"
)

const (
	BrokerAddress = "localhost:9092"
)

func NewWriterMessage(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(BrokerAddress),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
	}
}
