package people

import (
	"context"

	pb "github.com/Nickadimas79/kafka_testing/protobufs/people"
)

type Server struct {
	pb.UnimplementedPeopleServer
}

func (s *Server) CreatePerson(ctx context.Context, person *pb.Person) (*pb.Response, error) {
	return &pb.Response{
		RespMessage: "I WORKED!!!",
	}, nil
}
