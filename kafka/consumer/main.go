package consumer

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var testingTopic = "testing"

type Consumer struct {
	Con *kafka.Consumer
	SD  chan struct{}
}

// Instance singleton instance of a kafka consumer
var Instance Consumer

func New(conf kafka.ConfigMap) *Consumer {
	log.Println("creating Consumer")

	conf["group.id"] = "ontrack-assets-engine"
	conf["go.logs.channel.enable"] = true
	//conf["go.logs.channel"] = logKafka.LogChannel(cCtx)

	c, err := kafka.NewConsumer(&conf)
	if err != nil {
		log.Println("failed to create Consumer")
		os.Exit(1)
	}

	Instance = Consumer{
		Con: c,
		SD:  make(chan struct{}, 1),
	}

	return &Instance
}

func (c Consumer) TearDown() {
	// signal to shut down consumer loop
	c.SD <- struct{}{}
	log.Println("closing Consumer")

	// block to wait for loop to shut down
	<-c.SD

	err := c.Con.Close()
	if err != nil {
		log.Println("error closing Consumer")
	}
}

func (c Consumer) Consume() {
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
	// goroutine used to end loop when signal is sent
	go func() {
		_, ok := <-c.SD
		if ok {
			run = false
		}
	}()
	for run == true {
		ev := c.Con.Poll(100)
		if ev == nil {
			continue
		}

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
		}
	}

	c.SD <- struct{}{}
}
