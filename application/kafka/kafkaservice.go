package kafkaservice

import (
	"github.com/segmentio/kafka-go"
	"os"
)

func GetKafkaBrokerUrl() string {
	url, exists := os.LookupEnv("KAFKA_URL")
	if !exists {
		url = "localhost:9092"
	}
	return url
}

func NewWriterMessage(topic string) *kafka.Writer {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(GetKafkaBrokerUrl()),
		Topic:                  topic,
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
	return &writer
}

func NewReaderMessage(topic string, group string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{GetKafkaBrokerUrl()},
		Topic:   topic,
		GroupID: group,
	})
}
