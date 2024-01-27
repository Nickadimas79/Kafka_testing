package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())

	// start gRPC server on independent thread
	go grpc_server.StartGRPC()

	// this is a single instance of the Producer used for
	// producing all messages from this service
	pro := p.New(ccloud.ProducerConfig())

	// Consumer runs until it is told to stop (ctl-c) or errors
	con := c.New(ccloud.ConsumerConfig())
	go con.Consume(ctx)

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

			// clean up Producer object and connection
			go func() {
				defer wg.Done()

				log.Printf("flushing %v Messages/Requests from Producer\n", pro.Pro.Flush(1000))

				log.Println("closing Producer")
				pro.Pro.Close()
				log.Println("Producer closed")
			}()

			// shutdown of Consumer connection
			go func() {
				defer wg.Done()

				log.Println("canceling Consumer thread")
				cancel()

				// wait for consumer chan to signal finished shutting down
				<-con.ShutDown
				// time.Sleep(5 * time.Second)
			}()

			run = false
		}
	}

	wg.Wait()
	log.Println("system shut down")
}
