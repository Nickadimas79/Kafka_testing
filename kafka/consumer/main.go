package consumer

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nickadimas79/kafka_testing/kafka/ccloud"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Consumer() {
	fmt.Println("starting Consumer")

	wd, _ := os.Getwd()
	conf := ccloud.ReadConfig(wd + "/kafka/consumer/properties")
	conf["group.id"] = "kafka-testing"
	conf["auto.offset.reset"] = "earliest"

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "testing"
	err = c.SubscribeTopics([]string{topic}, nil)

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating Consumer\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}

			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}

	c.Close()
}
