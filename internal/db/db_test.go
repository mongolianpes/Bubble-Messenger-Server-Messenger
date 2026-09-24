package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestDSN() string {
	host := envOrDefault("TEST_DB_HOST", "localhost")
	port := envOrDefault("TEST_DB_PORT", "5432")
	user := envOrDefault("TEST_DB_USER", "postgres")
	password := envOrDefault("TEST_DB_PASSWORD", "123")
	dbname := envOrDefault("TEST_DB_NAME", "bubble")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func setupTestDB(t *testing.T) (*sql.DB, *PostgresStorage) {
	t.Helper()

	db, err := sql.Open("postgres", getTestDSN())
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, db.PingContext(ctx))

	_, err = db.ExecContext(ctx, `TRUNCATE TABLE messages RESTART IDENTITY CASCADE`)
	require.NoError(t, err)

	storage := &PostgresStorage{db: db}
	return db, storage
}

func TestSendMessage(t *testing.T) {
	db, storage := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	messageText := "text"
	err := storage.SendMessage(ctx, 1, 2, messageText)
	require.NoError(t, err)

	var sender, received int
	err = db.QueryRowContext(ctx,
		`SELECT sender_id, received_id FROM messages WHERE message = $1`, messageText,
	).Scan(&sender, &received)
	require.NoError(t, err)

	assert.Equal(t, 1, sender)
	assert.Equal(t, 2, received)
}

func TestCheckChats(t *testing.T) {
	db, storage := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	t.Run("have not messages", func(t *testing.T) {
		newMessages, err := storage.CheckChats(ctx, 1)
		require.NoError(t, err)

		require.Equal(t, len(newMessages), 0)
	})

	t.Run("have one messages", func(t *testing.T) {
		messageText := "text"
		_, err := db.ExecContext(ctx, sendMessage, 1, 2, messageText)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, sendMessage, 3, 2, messageText)
		require.NoError(t, err)

		newMessages, err := storage.CheckChats(ctx, 1)
		require.NoError(t, err)

		require.Equal(t, len(newMessages), 1)
		message := newMessages[0]
		require.Equal(t, message.SenderId, int64(1))
		require.Equal(t, message.ReceivedId, int64(2))
		require.Equal(t, message.Text, messageText)
	})

	t.Run("have many messages", func(t *testing.T) {
		messageText := "text"
		_, err := db.ExecContext(ctx, sendMessage, 1, 2, messageText)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, sendMessage, 3, 1, messageText)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, sendMessage, 3, 1, messageText)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, sendMessage, 3, 1, messageText)
		require.NoError(t, err)

		newMessages, err := storage.CheckChats(ctx, 1)
		require.NoError(t, err)

		require.Equal(t, len(newMessages), 5)
		message := newMessages[0]
		require.Equal(t, message.SenderId, int64(1))
		require.Equal(t, message.ReceivedId, int64(2))
		require.Equal(t, message.Text, messageText)

		message = newMessages[2]
		require.Equal(t, message.SenderId, int64(3))
		require.Equal(t, message.ReceivedId, int64(1))
		require.Equal(t, message.Text, messageText)
	})
}

func TestDelFewMessages(t *testing.T) {
	db, storage := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	t.Run("have not messages", func(t *testing.T) {
		err := storage.DelFewMessages(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("have one messages", func(t *testing.T) {
		messageText := "text"
		_, err := db.ExecContext(ctx, sendMessage, 1, 2, messageText)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, sendMessage, 3, 2, messageText)
		require.NoError(t, err)

		err = storage.DelFewMessages(ctx, 1)
		require.NoError(t, err)

		var receivedID int
		err = db.QueryRowContext(ctx, "SELECT received_id FROM messages WHERE sender_id = $1", 3).Scan(&receivedID)
		require.NoError(t, err)
		require.Equal(t, receivedID, 2)

		err = db.QueryRowContext(ctx, "SELECT received_id FROM messages WHERE sender_id = $1", 1).Scan(&receivedID)
		require.Error(t, err)
		require.Equal(t, err, sql.ErrNoRows)
	})

	t.Run("have many messages", func(t *testing.T) {
		messageText := "text"
		_, err := db.ExecContext(ctx, sendMessage, 1, 2, messageText)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, sendMessage, 3, 1, messageText)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, sendMessage, 3, 1, messageText)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, sendMessage, 3, 1, messageText)
		require.NoError(t, err)

		err = storage.DelFewMessages(ctx, 1)
		require.NoError(t, err)
	})
}
