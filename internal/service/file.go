package service

import (
	"context"
	"fmt"
	pb "messenger/proto"
)

const designationSendFile = "p"

func (s *MessengerService) SendFile(ctx context.Context, req *pb.SendFileRequest) (*pb.SendFileResponse, error) {
	textMessage := fmt.Sprintf("%s\\%s\\%s", designationSendFile, req.FileName, req.DatabasePath)
	if err := s.SQLStorage.SendMessage(ctx, int(req.SenderId), int(req.ReceivedId), textMessage); err != nil {
		return nil, err
	}

	return &pb.SendFileResponse{}, nil
}
