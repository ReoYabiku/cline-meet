package repository

import (
	"context"

	"github.com/cline-meet/backend/internal/domain/model"
	"github.com/google/uuid"
)

// MessageRepository defines the interface for Redis-based message operations
type Message interface {
	// SaveChatMessage saves a chat message to Redis List
	// Messages are stored in Redis with automatic expiration
	SaveChatMessage(ctx context.Context, message *model.Message) error
	
	// GetChatHistory retrieves recent chat history for a room from Redis
	// Returns the most recent messages (up to limit) in chronological order
	GetChatHistory(ctx context.Context, roomID uuid.UUID, limit int) ([]*model.Message, error)
	
	// DeleteChatHistory deletes all chat history for a room from Redis
	DeleteChatHistory(ctx context.Context, roomID uuid.UUID) error
}
