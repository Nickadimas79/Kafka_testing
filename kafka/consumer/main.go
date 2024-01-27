package consumer

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var testingTopic = "testing"

type Consumer struct {
	Con      *kafka.Consumer
	ShutDown chan struct{}
}

// Instance singleton instance of a kafka consumer
var Instance Consumer

func New(conf kafka.ConfigMap) *Consumer {
	log.Println("creating Consumer")

	conf["group.id"] = "ontrack-assets-engine"
	conf["go.logs.channel.enable"] = true
	// conf["go.logs.channel"] = logKafka.LogChannel(cCtx)

	c, err := kafka.NewConsumer(&conf)
	if err != nil {
		log.Println("failed to create Consumer")
		os.Exit(1)
	}

	Instance = Consumer{
		Con:      c,
		ShutDown: make(chan struct{}, 1),
	}

	return &Instance
}

func (c Consumer) Consume(ctx context.Context) {
	log.Println("starting Consumer")
	topics := []string{
		testingTopic,
	}

	err := c.Con.SubscribeTopics(topics, nil)
	if err != nil {
		log.Println("failed to subscribe to Topics")
		os.Exit(1)
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run == true {
		ev := c.Con.Poll(100)
		if ev == nil {
			continue
		}
		// select for closing down context
		select {
		case <-ctx.Done():
			log.Println("Consumer thread canceled")
			run = false
		default:
			// switch to check event type
			switch msg := ev.(type) {
			case *kafka.Message:
				// switch to check topic type
				switch *msg.TopicPartition.Topic {
				case testingTopic:
					log.Printf("consumed Message from Topic %s\n", *msg.TopicPartition.Topic)
					log.Printf("with value: %+v\n", string(msg.Value))
				}
			case kafka.Error:
				log.Println("Kafka error:", msg)
			}
		}
	}

	log.Println("closing Consumer")
	err = c.Con.Close()
	if err != nil {
		log.Println("error closing Consumer")
	}

	log.Println("Consumer closed")
	c.ShutDown <- struct{}{}
}
