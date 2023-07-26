proto:
		protoc --go_out=protobufs --go-grpc_out=protobufs protos/*.proto

run:
		go run cmd/main.go