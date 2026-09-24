package db

import (
	"context"
	"time"

	pb "messenger/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	sendMessage = "INSERT INTO messages (sender_id, received_id, message) VALUES ($1, $2, $3)"
	checkChats  = `SELECT sender_id, received_id, create_at, message
		FROM messages
		WHERE sender_id = $1 OR received_id = $1
		ORDER BY create_at ASC
		LIMIT 5`
	delMessages = `DELETE FROM messages
		WHERE id IN (
			SELECT id
			FROM messages
			WHERE sender_id = $1 OR received_id = $1
			ORDER BY create_at ASC
			LIMIT 5
		)`
)

func (s *PostgresStorage) SendMessage(ctx context.Context, senderID, receivedID int, text string) error {
	_, err := s.db.ExecContext(ctx, sendMessage, senderID, receivedID, text)
	return err
}

func (s *PostgresStorage) CheckChats(ctx context.Context, userID int) ([]*pb.NewMessage, error) {
	rows, err := s.db.QueryContext(ctx, checkChats, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newMessages := []*pb.NewMessage{}
	senderID := 0
	receivedID := 0
	var createAt time.Time
	text := ""
	for rows.Next() {
		if err := rows.Scan(&senderID, &receivedID, &createAt, &text); err != nil {
			return nil, err
		}

		newMessages = append(newMessages, &pb.NewMessage{
			Text:       text,
			SenderId:   int64(senderID),
			ReceivedId: int64(receivedID),
			CreateAt:   timestamppb.New(createAt),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return newMessages, nil
}

func (s *PostgresStorage) DelFewMessages(ctx context.Context, userID int) error {
	_, err := s.db.ExecContext(ctx, delMessages, userID)
	return err
}
