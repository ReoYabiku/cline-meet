package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cline-meet/backend/internal/domain/model"
	"github.com/cline-meet/backend/internal/domain/repository"
	"github.com/cline-meet/backend/internal/domain/service"
	"github.com/google/uuid"
)

// Room handles room-related business logic
type Room struct {
	roomRepo         repository.Room
	userRepo         repository.User
	realtimeNotifier service.RealtimeNotifier
	sessionManager   service.SessionManager
}

// NewRoom creates a new Room usecase
func NewRoom(
	roomRepo repository.Room,
	userRepo repository.User,
	realtimeNotifier service.RealtimeNotifier,
	sessionManager service.SessionManager,
) *Room {
	return &Room{
		roomRepo:         roomRepo,
		userRepo:         userRepo,
		realtimeNotifier: realtimeNotifier,
		sessionManager:   sessionManager,
	}
}

// CreateRoom creates a new meeting room
func (r *Room) CreateRoom(ctx context.Context, hostID uuid.UUID, name string, isWaitingRoom bool) (*model.Room, error) {
	// Validate host exists
	_, err := r.userRepo.GetByID(ctx, hostID)
	if err != nil {
		return nil, fmt.Errorf("host not found: %w", err)
	}

	// Create room
	room := model.NewRoom(name, hostID, isWaitingRoom)

	// Add host as first participant
	if err := room.AddParticipant(hostID); err != nil {
		return nil, fmt.Errorf("failed to add host to room: %w", err)
	}

	// Save to repository
	if err := r.roomRepo.Create(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	return room, nil
}

// JoinRoom adds a user to a room
func (r *Room) JoinRoom(ctx context.Context, userID, roomID uuid.UUID) error {
	// Get user
	user, err := r.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Get room
	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	// Check if room is expired
	if room.IsExpired() {
		return errors.New("room has expired")
	}

	// Add participant to room
	if err := room.AddParticipant(userID); err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}

	// Update room in repository
	if err := r.roomRepo.Update(ctx, room); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	// Create or update user session
	session := &service.UserSession{
		UserID:   userID,
		RoomID:   roomID,
		IsHost:   room.IsHost(userID),
		IsMuted:  false,
		LastSeen: time.Now().Unix(),
	}

	if err := r.sessionManager.UpdateSession(ctx, session); err != nil {
		// Log error but don't fail the join operation
		// Session management is not critical for basic functionality
	}

	// Notify other participants
	if err := r.realtimeNotifier.NotifyRoomJoined(ctx, roomID, userID, user.Name); err != nil {
		// Log error but don't fail the join operation
		// Real-time notification failure shouldn't prevent joining
	}

	return nil
}

// LeaveRoom removes a user from a room
func (r *Room) LeaveRoom(ctx context.Context, userID, roomID uuid.UUID) error {
	// Get user
	user, err := r.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Get room
	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	// Remove participant from room
	if err := room.RemoveParticipant(userID); err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	// If room is empty, delete it
	if room.GetParticipantCount() == 0 {
		if err := r.roomRepo.Delete(ctx, roomID); err != nil {
			return fmt.Errorf("failed to delete empty room: %w", err)
		}
	} else {
		// Update room in repository
		if err := r.roomRepo.Update(ctx, room); err != nil {
			return fmt.Errorf("failed to update room: %w", err)
		}
	}

	// Delete user session
	if err := r.sessionManager.DeleteSession(ctx, userID); err != nil {
		// Log error but don't fail the leave operation
	}

	// Notify other participants
	if err := r.realtimeNotifier.NotifyRoomLeft(ctx, roomID, userID, user.Name); err != nil {
		// Log error but don't fail the leave operation
	}

	return nil
}

// MuteParticipant mutes a participant (only host can do this)
func (r *Room) MuteParticipant(ctx context.Context, hostID, roomID, targetUserID uuid.UUID) error {
	// Get room
	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	// Mute participant
	if err := room.MuteParticipant(hostID, targetUserID); err != nil {
		return fmt.Errorf("failed to mute participant: %w", err)
	}

	// Update room in repository
	if err := r.roomRepo.Update(ctx, room); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	// Update session
	session, err := r.sessionManager.GetSession(ctx, targetUserID)
	if err == nil {
		session.IsMuted = true
		r.sessionManager.UpdateSession(ctx, session)
	}

	// Notify participants
	if err := r.realtimeNotifier.NotifyUserMuted(ctx, roomID, targetUserID, true); err != nil {
		// Log error but don't fail the operation
	}

	return nil
}

// UnmuteParticipant unmutes a participant
func (r *Room) UnmuteParticipant(ctx context.Context, userID, roomID uuid.UUID) error {
	// Get room
	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	// Unmute participant
	if err := room.UnmuteParticipant(userID); err != nil {
		return fmt.Errorf("failed to unmute participant: %w", err)
	}

	// Update room in repository
	if err := r.roomRepo.Update(ctx, room); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	// Update session
	session, err := r.sessionManager.GetSession(ctx, userID)
	if err == nil {
		session.IsMuted = false
		r.sessionManager.UpdateSession(ctx, session)
	}

	// Notify participants
	if err := r.realtimeNotifier.NotifyUserMuted(ctx, roomID, userID, false); err != nil {
		// Log error but don't fail the operation
	}

	return nil
}

// GetRoom retrieves a room by ID
func (r *Room) GetRoom(ctx context.Context, roomID uuid.UUID) (*model.Room, error) {
	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	// Check if room is expired
	if room.IsExpired() {
		return nil, errors.New("room has expired")
	}

	return room, nil
}

// GetUserRooms retrieves all rooms where the user is the host
func (r *Room) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]*model.Room, error) {
	rooms, err := r.roomRepo.GetByHostID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rooms: %w", err)
	}

	// Filter out expired rooms
	var activeRooms []*model.Room
	for _, room := range rooms {
		if !room.IsExpired() {
			activeRooms = append(activeRooms, room)
		}
	}

	return activeRooms, nil
}

// ExtendRoomExpiry extends the expiry time of a room (only host can do this)
func (r *Room) ExtendRoomExpiry(ctx context.Context, hostID, roomID uuid.UUID, hours int) error {
	// Get room
	room, err := r.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}

	// Check if user is host
	if !room.IsHost(hostID) {
		return errors.New("only host can extend room expiry")
	}

	// Extend expiry
	room.ExtendExpiry(time.Duration(hours) * time.Hour)

	// Update room in repository
	if err := r.roomRepo.Update(ctx, room); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	// Notify participants about room update
	if err := r.realtimeNotifier.NotifyRoomUpdate(ctx, room); err != nil {
		// Log error but don't fail the operation
	}

	return nil
}
