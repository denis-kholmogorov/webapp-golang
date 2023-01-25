package kafkasender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"web/application/domain"
)

const (
	topic          = "message-log"
	broker1Address = "localhost:9092"
)

var sender KafkaSender
var isInitialized bool

type KafkaSender struct {
	kafkaWriter *kafka.Writer
}

func NewKafkaSender() *KafkaSender {
	mt := sync.Mutex{}
	mt.Lock()
	if !isInitialized {
		kafkaWriter := kafka.Writer{
			Addr:         kafka.TCP(broker1Address),
			Topic:        topic,
			RequiredAcks: kafka.RequireAll,
		}
		sender = KafkaSender{kafkaWriter: &kafkaWriter}
		isInitialized = true
	}
	mt.Unlock()
	return &sender
}

func (k *KafkaSender) SendMessage(messages ...any) {
	for _, m := range messages {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		marshal, _ := json.Marshal(m)
		err := k.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(uuid.NewString()),
			Value: marshal,
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}
		fmt.Println("writes:", string(marshal))
	}
}

func Consume1(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value

		person := domain.Account{}
		err = json.Unmarshal(msg.Value, &person)
		if err != nil {
			log.Println("Kafka can't read message")
		}

		fmt.Println("received consumer 1: ", person)
		fmt.Println("received consumer 1 message : {", "msg.Partition", msg.Partition, "msg.Key", string(msg.Key), "msg.Offset:", msg.Offset, "msg.HighWaterMark:", msg.HighWaterMark, "msg.Topic:", msg.Topic, "msg.Headers", msg.Headers, "}")
	}
}

func Consume2(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value

		person := domain.Account{}
		err = json.Unmarshal(msg.Value, &person)
		if err != nil {
			log.Println("Kafka can't read message")
		}
		fmt.Println("received consumer 2: ", person)
		fmt.Println("received consumer 2 message : {", "msg.Partition", msg.Partition, "msg.Key", string(msg.Key), "msg.Offset:", msg.Offset, "msg.HighWaterMark:", msg.HighWaterMark, "msg.Topic:", msg.Topic, "msg.Headers", msg.Headers, "}")
	}
}
