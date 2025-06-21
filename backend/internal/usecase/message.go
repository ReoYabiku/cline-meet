package usecase

import (
	"context"
	"fmt"

	"github.com/cline-meet/backend/internal/domain/model"
	"github.com/cline-meet/backend/internal/domain/repository"
	"github.com/cline-meet/backend/internal/domain/service"
	"github.com/google/uuid"
)

// Message handles message-related business logic
type Message struct {
	messageRepo      repository.Message
	roomRepo         repository.Room
	userRepo         repository.User
	realtimeNotifier service.RealtimeNotifier
}

// NewMessage creates a new Message usecase
func NewMessage(
	messageRepo repository.Message,
	roomRepo repository.Room,
	userRepo repository.User,
	realtimeNotifier service.RealtimeNotifier,
) *Message {
	return &Message{
		messageRepo:      messageRepo,
		roomRepo:         roomRepo,
		userRepo:         userRepo,
		realtimeNotifier: realtimeNotifier,
	}
}

// SendChatMessage sends a chat message to a room
func (c *Message) SendMessage(ctx context.Context, senderID, roomID uuid.UUID, messageText string) error {
	// Get user
	user, err := c.userRepo.GetByID(ctx, senderID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Get room
	room, err := c.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	// Check if room is expired
	if room.IsExpired() {
		return fmt.Errorf("room has expired")
	}

	// Check if user is a participant
	if !room.IsParticipant(senderID) {
		return fmt.Errorf("user is not a participant in this room")
	}

	// Create chat message
	message := model.NewChatMessage(senderID, roomID, messageText, user.Name)

	// Validate message
	if !message.IsValid() {
		return fmt.Errorf("invalid message")
	}

	// Save message to Redis
	if err := c.messageRepo.SaveChatMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	// Broadcast message to all room participants
	if err := c.realtimeNotifier.BroadcastChatMessage(ctx, message); err != nil {
		// Log error but don't fail the send operation
		// Real-time notification failure shouldn't prevent message saving
	}

	return nil
}

// GetChatHistory retrieves chat history for a room
func (c *Message) GetHistory(ctx context.Context, userID, roomID uuid.UUID, limit int) ([]*model.Message, error) {
	// Get room
	room, err := c.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	// Check if user is a participant
	if !room.IsParticipant(userID) {
		return nil, fmt.Errorf("user is not a participant in this room")
	}

	// Set default limit if not specified
	if limit <= 0 {
		limit = 50 // Default to last 50 messages
	}

	// Get chat history from Redis
	messages, err := c.messageRepo.GetChatHistory(ctx, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}

	return messages, nil
}

// DeleteChatHistory deletes all chat history for a room (only host can do this)
func (c *Message) DeleteHistory(ctx context.Context, hostID, roomID uuid.UUID) error {
	// Get room
	room, err := c.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	// Check if user is the host
	if !room.IsHost(hostID) {
		return fmt.Errorf("only host can delete chat history")
	}

	// Delete chat history from Redis
	if err := c.messageRepo.DeleteChatHistory(ctx, roomID); err != nil {
		return fmt.Errorf("failed to delete chat history: %w", err)
	}

	return nil
}
