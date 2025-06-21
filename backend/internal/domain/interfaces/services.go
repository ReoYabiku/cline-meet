package interfaces

import (
	"context"

	"github.com/cline-meet/backend/internal/domain/model"
	"github.com/google/uuid"
)

// RealtimeNotifier defines the interface for real-time notifications
// This abstracts the WebSocket communication from the application layer
type RealtimeNotifier interface {
	// NotifyRoomJoined notifies all participants that a user joined the room
	NotifyRoomJoined(ctx context.Context, roomID, userID uuid.UUID, userName string) error
	
	// NotifyRoomLeft notifies all participants that a user left the room
	NotifyRoomLeft(ctx context.Context, roomID, userID uuid.UUID, userName string) error
	
	// NotifyUserMuted notifies all participants that a user was muted
	NotifyUserMuted(ctx context.Context, roomID, userID uuid.UUID, isMuted bool) error
	
	// BroadcastChatMessage broadcasts a chat message to all room participants
	BroadcastChatMessage(ctx context.Context, message *model.Message) error
	
	// SendDirectMessage sends a direct message to a specific user (for WebRTC signaling)
	SendDirectMessage(ctx context.Context, message *model.Message) error
	
	// NotifyRoomUpdate notifies participants about room setting changes
	NotifyRoomUpdate(ctx context.Context, room *model.Room) error
}

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
