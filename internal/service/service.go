package service

import (
	"messenger/internal/db"
	pb "messenger/proto"
)

type MessengerService struct {
	pb.UnimplementedMessengerServiceServer
	SQLStorage db.SQLStorage
}

func NewMessengerService(messagesStorage db.SQLStorage) *MessengerService {
	return &MessengerService{
		SQLStorage: messagesStorage,
	}
}
