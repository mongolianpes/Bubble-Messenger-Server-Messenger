package db

import (
	"context"
	"database/sql"
	"fmt"
	pb "messenger/proto"
	"os"
)

type PostgresStorage struct {
	db *sql.DB
}

type SQLStorage interface {
	SendMessage(ctx context.Context, senderID, receivedID int, text string) error
	CheckChats(ctx context.Context, userID int) ([]*pb.NewMessage, error)
	DelFewMessages(ctx context.Context, userID int) error
}

func NewPostgresStorage() (*PostgresStorage, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{
		db: db,
	}, nil
}

func connectToDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
