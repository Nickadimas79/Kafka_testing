package main

import (
	"fmt"
	c "github.com/Nickadimas79/kafka_testing/kafka/consumer"
	"log"
	"net"
	"os"

	"github.com/Nickadimas79/kafka_testing/pkg/people"

	p "github.com/Nickadimas79/kafka_testing/kafka/producer"
	peoplePB "github.com/Nickadimas79/kafka_testing/protobufs/people"

	"google.golang.org/grpc"
)

const gRPCPort = ":8080"

func main() {
	fmt.Println("::STARTING PRODUCER/CONSUMER APP::")

	// this is a single run of the Producer
	go p.Producer()
	// Consumer runs until it is told to stop (ctl-c) or errors
	go c.Consumer()

	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Println("failed to start TCP listener")
		os.Exit(1)
	}

	peopleServer := people.Server{}

	grpcServer := grpc.NewServer()

	peoplePB.RegisterPeopleServer(grpcServer, &peopleServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Println("failed to start gRPC server")
		os.Exit(1)
	}
}
