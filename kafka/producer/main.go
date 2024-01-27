package producer

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	Pro *kafka.Producer
}

// Instance singleton instance of a kafka producer
var Instance Producer

func New(configMap kafka.ConfigMap) *Producer {
	log.Println("creating Producer")

	p, err := kafka.NewProducer(&configMap)

	if err != nil {
		log.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	Instance = Producer{
		Pro: p,
	}

	return &Instance
}

func (p Producer) Produce(ctx context.Context, data *kafka.Message) {
	errChan := make(chan kafka.Event)

	err := p.Pro.Produce(data, errChan)
	if err != nil {
		log.Println("Kafka Message not enqueued/produced:", err)
	}

	e := <-errChan

	switch ev := e.(type) {
	case kafka.Error:
		log.Println("Producer error")
		close(errChan)

	case *kafka.Message:
		log.Printf("produced Message to Topic %s\n", *ev.TopicPartition.Topic)
		log.Printf("with key = %s offset = %v value = %v\n",
			string(ev.Key), ev.TopicPartition.Offset, string(ev.Value))

		close(errChan)
	}
}

// BuildMsg builds a Kafka msg from any data type passed for the topic you pass.
// Can't be bound to any interfaces or objects, IE as a method, because of the
// nature of generics at this time. (?)
func BuildMsg[T any](data *T, topic string) *kafka.Message {
	d, _ := json.Marshal(data)

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte{},
		Value: d,
	}
}
