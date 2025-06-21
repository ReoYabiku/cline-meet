package repository

import (
	"context"

	"github.com/cline-meet/backend/internal/domain/model"
	"github.com/google/uuid"
)

// RoomRepository defines the interface for room data operations
type Room interface {
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
