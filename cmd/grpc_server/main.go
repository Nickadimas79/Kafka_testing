package grpc_server

import (
	"log"
	"net"
	"os"

	"github.com/Nickadimas79/kafka_testing/pkg/people"
	peoplePB "github.com/Nickadimas79/kafka_testing/protobufs/people"

	"google.golang.org/grpc"
)

const gRPCPort = ":8080"

func StartGRPC() {
	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Println("failed to start TCP listener")
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	peopleServ := people.Server{}

	peoplePB.RegisterPeopleServer(grpcServer, &peopleServ)

	if err := grpcServer.Serve(lis); err != nil {
		log.Println("failed to start gRPC server")
		os.Exit(1)
	}
}
