package producer

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	Instance = Producer{
		Pro: p,
	}

	return &Instance
}

func (p Producer) TearDown() {
	log.Printf("flushing %v Messages/Requests from Producer\n", p.Pro.Flush(10000))

	log.Println("closing Producer")
	p.Pro.Close()
}

func (p Producer) Produce(data *kafka.Message) {
	errChan := make(chan kafka.Event)

	err := p.Pro.Produce(data, errChan)
	if err != nil {
		log.Println("Kafka Message not enqueued")
		//return err
	}

	e := <-errChan

	switch ev := e.(type) {
	case kafka.Error:
		log.Println("Producer error")
		close(errChan)
		//return ev
	case *kafka.Message:
		log.Printf("produced Message to Topic %s\n", *ev.TopicPartition.Topic)
		log.Printf("with key = %s offset = %v value = %v\n",
			string(ev.Key), ev.TopicPartition.Offset, string(ev.Value))

		close(errChan)
	}

	//return nil
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
