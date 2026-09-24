package service

import (
	"context"
	pb "messenger/proto"
)

func (s *MessengerService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	if err := s.SQLStorage.SendMessage(ctx, int(req.SenderId), int(req.ReceivedId), req.Text); err != nil {
		return nil, err
	}

	return &pb.SendMessageResponse{}, nil
}

func (s *MessengerService) CheckChats(ctx context.Context, req *pb.CheckChatsRequest) (*pb.CheckChatsResponse, error) {
	newMessages, err := s.SQLStorage.CheckChats(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.CheckChatsResponse{
		Messages: newMessages,
	}, nil
}

func (s *MessengerService) DelMessage(ctx context.Context, req *pb.DelMessageRequest) (*pb.DelMessageResponse, error) {
	if err := s.SQLStorage.DelFewMessages(ctx, int(req.UserId)); err != nil {
		return nil, err
	}

	return &pb.DelMessageResponse{}, nil
}
