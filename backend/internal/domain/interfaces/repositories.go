package interfaces

import (
	"context"

	"github.com/cline-meet/backend/internal/domain/model"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *model.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	
	// GetByGoogleID retrieves a user by Google ID
	GetByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *model.User) error
	
	// Delete deletes a user
	Delete(ctx context.Context, id uuid.UUID) error
}

// RoomRepository defines the interface for room data operations
type RoomRepository interface {
	// Create creates a new room
	Create(ctx context.Context, room *model.Room) error
	
	// GetByID retrieves a room by ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error)
	
	// GetByHostID retrieves rooms by host ID
	GetByHostID(ctx context.Context, hostID uuid.UUID) ([]*model.Room, error)
	
	// Update updates an existing room
	Update(ctx context.Context, room *model.Room) error
	
	// Delete deletes a room
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetActiveRooms retrieves all active (non-expired) rooms
	GetActiveRooms(ctx context.Context) ([]*model.Room, error)
	
	// CleanupExpiredRooms removes expired rooms
	// This method is called by a background scheduler every 1 hour
	CleanupExpiredRooms(ctx context.Context) error
}

// MessageRepository defines the interface for Redis-based message operations
type MessageRepository interface {
	// SaveChatMessage saves a chat message to Redis List
	// Messages are stored in Redis with automatic expiration
	SaveChatMessage(ctx context.Context, message *model.Message) error
	
	// GetChatHistory retrieves recent chat history for a room from Redis
	// Returns the most recent messages (up to limit) in chronological order
	GetChatHistory(ctx context.Context, roomID uuid.UUID, limit int) ([]*model.Message, error)
	
	// DeleteChatHistory deletes all chat history for a room from Redis
	DeleteChatHistory(ctx context.Context, roomID uuid.UUID) error
}
