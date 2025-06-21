package service

import (
	"context"

	"github.com/google/uuid"
)

// SessionManager defines the interface for managing user sessions
type SessionManager interface {
	// CreateSession creates a new user session
	CreateSession(ctx context.Context, userID uuid.UUID, connectionID string) error
	
	// GetSession retrieves a user session
	GetSession(ctx context.Context, userID uuid.UUID) (*UserSession, error)
	
	// UpdateSession updates a user session
	UpdateSession(ctx context.Context, session *UserSession) error
	
	// DeleteSession deletes a user session
	DeleteSession(ctx context.Context, userID uuid.UUID) error
	
	// GetActiveUsers returns all active users in a room
	GetActiveUsers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error)
}

// UserSession represents an active user session
type UserSession struct {
	UserID       uuid.UUID `json:"userId"`
	ConnectionID string    `json:"connectionId"`
	RoomID       uuid.UUID `json:"roomId,omitempty"`
	IsHost       bool      `json:"isHost"`
	IsMuted      bool      `json:"isMuted"`
	ServerPod    string    `json:"serverPod"`
	LastSeen     int64     `json:"lastSeen"` // Unix timestamp
}
