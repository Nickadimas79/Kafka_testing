package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Nickadimas79/kafka_testing/cmd/grpc_server"
	"github.com/Nickadimas79/kafka_testing/kafka/ccloud"
	c "github.com/Nickadimas79/kafka_testing/kafka/consumer"
	p "github.com/Nickadimas79/kafka_testing/kafka/producer"
)

func main() {
	fmt.Println("::STARTING APP::")

	// start gRPC server on independent thread
	go grpc_server.StartGRPC()

	// this is a single instance of the Producer used for
	// producing all messages from this service
	pro := p.New(ccloud.ProducerConfig())

	// Consumer runs until it is told to stop (ctl-c) or errors
	con := c.New(ccloud.ConsumerConfig())
	go con.Consume()

	// Set up a channel for handling Ctrl-C, etc...
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	run := true

	// loop checking for OS signals
	for run == true {
		select {
		// graceful shutdown of Producer/Consumer
		case sig := <-sigChan:
			wg.Add(2)
			log.Printf("caught signal %v\n", sig)
			go func() {
				defer wg.Done()
				pro.TearDown()
			}()
			go func() {
				defer wg.Done()
				con.TearDown()
			}()

			run = false
		}
	}

	wg.Wait()
	log.Println("system shut down")
}
