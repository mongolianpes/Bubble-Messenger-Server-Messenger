package main

import (
	"messenger/internal/db"
	"messenger/internal/service"
	"net"

	pb "messenger/proto"

	"google.golang.org/grpc"
)

func main() {
	messagesStorage, err := db.NewPostgresStorage()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessengerServiceServer(grpcServer, service.NewMessengerService(messagesStorage))

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
