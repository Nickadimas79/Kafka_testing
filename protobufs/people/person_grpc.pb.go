// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: protos/person.proto

package people

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PeopleClient is the client API for People service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PeopleClient interface {
	//  rpc GetPerson(PersonReq) returns (Person) {}
	CreatePerson(ctx context.Context, in *Person, opts ...grpc.CallOption) (*Response, error)
}

type peopleClient struct {
	cc grpc.ClientConnInterface
}

func NewPeopleClient(cc grpc.ClientConnInterface) PeopleClient {
	return &peopleClient{cc}
}

func (c *peopleClient) CreatePerson(ctx context.Context, in *Person, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/kafka_testing.People/CreatePerson", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PeopleServer is the server API for People service.
// All implementations must embed UnimplementedPeopleServer
// for forward compatibility
type PeopleServer interface {
	//  rpc GetPerson(PersonReq) returns (Person) {}
	CreatePerson(context.Context, *Person) (*Response, error)
	mustEmbedUnimplementedPeopleServer()
}

// UnimplementedPeopleServer must be embedded to have forward compatible implementations.
type UnimplementedPeopleServer struct {
}

func (UnimplementedPeopleServer) CreatePerson(context.Context, *Person) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePerson not implemented")
}
func (UnimplementedPeopleServer) mustEmbedUnimplementedPeopleServer() {}

// UnsafePeopleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PeopleServer will
// result in compilation errors.
type UnsafePeopleServer interface {
	mustEmbedUnimplementedPeopleServer()
}

func RegisterPeopleServer(s grpc.ServiceRegistrar, srv PeopleServer) {
	s.RegisterService(&People_ServiceDesc, srv)
}

func _People_CreatePerson_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Person)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeopleServer).CreatePerson(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kafka_testing.People/CreatePerson",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeopleServer).CreatePerson(ctx, req.(*Person))
	}
	return interceptor(ctx, in, info, handler)
}

// People_ServiceDesc is the grpc.ServiceDesc for People service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var People_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kafka_testing.People",
	HandlerType: (*PeopleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePerson",
			Handler:    _People_CreatePerson_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/person.proto",
}
