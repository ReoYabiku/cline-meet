package service

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
