package main

import (
	"log/slog"
	"messenger/internal/db"
	"messenger/internal/service"
	"net"

	pb "messenger/proto"

	"google.golang.org/grpc"
)

func main() {
	messagesStorage, err := db.NewPostgresStorage()
	if err != nil {
		slog.Error("Connect to DB", "error", err)
		return
	}

	result, err := messagesStorage.DeleteOldMessages()
	if err != nil {
		slog.Error("Delete old messages", "error", err)
	} else {
		slog.Info("Delete old messages", "rowsAffected", result)
	}

	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		slog.Error("Start listen port", "error", err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessengerServiceServer(grpcServer, service.NewMessengerService(messagesStorage))

	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("Serve gRPC service", "error", err)
	}
}
