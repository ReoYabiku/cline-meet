package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Room represents a meeting room
type Room struct {
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	HostID        uuid.UUID     `json:"hostId"`
	IsWaitingRoom bool          `json:"isWaitingRoom"`
	CreatedAt     time.Time     `json:"createdAt"`
	ExpiresAt     time.Time     `json:"expiresAt"`
	Participants  []Participant `json:"participants"`
	MaxCapacity   int           `json:"maxCapacity"`
}

// Participant represents a participant in a room
type Participant struct {
	UserID   uuid.UUID `json:"userId"`
	IsHost   bool      `json:"isHost"`
	IsMuted  bool      `json:"isMuted"`
	JoinedAt time.Time `json:"joinedAt"`
}

// NewRoom creates a new room
func NewRoom(name string, hostID uuid.UUID, isWaitingRoom bool) *Room {
	now := time.Now()
	return &Room{
		ID:            uuid.New(),
		Name:          name,
		HostID:        hostID,
		IsWaitingRoom: isWaitingRoom,
		CreatedAt:     now,
		ExpiresAt:     now.Add(24 * time.Hour), // 24時間後に期限切れ
		Participants:  []Participant{},
		MaxCapacity:   10, // Google Meetクローンの要件
	}
}

// AddParticipant adds a participant to the room
func (r *Room) AddParticipant(userID uuid.UUID) error {
	// 既に参加しているかチェック
	for _, p := range r.Participants {
		if p.UserID == userID {
			return errors.New("user already in room")
		}
	}

	// 定員チェック
	if len(r.Participants) >= r.MaxCapacity {
		return errors.New("room is full")
	}

	// 期限チェック
	if time.Now().After(r.ExpiresAt) {
		return errors.New("room has expired")
	}

	participant := Participant{
		UserID:   userID,
		IsHost:   userID == r.HostID,
		IsMuted:  false,
		JoinedAt: time.Now(),
	}

	r.Participants = append(r.Participants, participant)
	return nil
}

// RemoveParticipant removes a participant from the room
func (r *Room) RemoveParticipant(userID uuid.UUID) error {
	for i, p := range r.Participants {
		if p.UserID == userID {
			// スライスから削除
			r.Participants = append(r.Participants[:i], r.Participants[i+1:]...)
			return nil
		}
	}
	return errors.New("participant not found")
}

// GetParticipant returns a participant by user ID
func (r *Room) GetParticipant(userID uuid.UUID) (*Participant, error) {
	for i, p := range r.Participants {
		if p.UserID == userID {
			return &r.Participants[i], nil
		}
	}
	return nil, errors.New("participant not found")
}

// MuteParticipant mutes a participant (only host can do this)
func (r *Room) MuteParticipant(hostID, targetUserID uuid.UUID) error {
	if hostID != r.HostID {
		return errors.New("only host can mute participants")
	}

	participant, err := r.GetParticipant(targetUserID)
	if err != nil {
		return err
	}

	participant.IsMuted = true
	return nil
}

// UnmuteParticipant unmutes a participant
func (r *Room) UnmuteParticipant(userID uuid.UUID) error {
	participant, err := r.GetParticipant(userID)
	if err != nil {
		return err
	}

	participant.IsMuted = false
	return nil
}

// IsHost checks if a user is the host
func (r *Room) IsHost(userID uuid.UUID) bool {
	return r.HostID == userID
}

// IsParticipant checks if a user is a participant
func (r *Room) IsParticipant(userID uuid.UUID) bool {
	_, err := r.GetParticipant(userID)
	return err == nil
}

// GetParticipantCount returns the number of participants
func (r *Room) GetParticipantCount() int {
	return len(r.Participants)
}

// IsExpired checks if the room has expired
func (r *Room) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsFull checks if the room is at capacity
func (r *Room) IsFull() bool {
	return len(r.Participants) >= r.MaxCapacity
}

// ExtendExpiry extends the room expiry time
func (r *Room) ExtendExpiry(duration time.Duration) {
	r.ExpiresAt = r.ExpiresAt.Add(duration)
}
