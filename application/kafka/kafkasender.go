package kafkasender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	topic          = "message-log"
	broker1Address = "localhost:9092"
)

var kafkaStruct KafkaSender
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
		kafkaStruct = KafkaSender{kafkaWriter: &kafkaWriter}
		isInitialized = true
	}
	mt.Unlock()
	return &kafkaStruct
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

		// log a confirmation once the message is written
		fmt.Println("writes:", marshal)
	}
}

func Consume(ctx context.Context) {
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
		fmt.Println("received: ", string(msg.Value))
	}
}
