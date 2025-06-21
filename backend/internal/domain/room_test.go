package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewRoom(t *testing.T) {
	name := "Test Room"
	hostID := uuid.New()
	isWaitingRoom := true

	room := NewRoom(name, hostID, isWaitingRoom)

	if room.Name != name {
		t.Errorf("Expected Name %s, got %s", name, room.Name)
	}
	if room.HostID != hostID {
		t.Errorf("Expected HostID %s, got %s", hostID, room.HostID)
	}
	if room.IsWaitingRoom != isWaitingRoom {
		t.Errorf("Expected IsWaitingRoom %v, got %v", isWaitingRoom, room.IsWaitingRoom)
	}
	if room.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}
	if room.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if room.ExpiresAt.IsZero() {
		t.Error("Expected ExpiresAt to be set")
	}
	if room.MaxCapacity != 10 {
		t.Errorf("Expected MaxCapacity 10, got %d", room.MaxCapacity)
	}
	if len(room.Participants) != 0 {
		t.Errorf("Expected empty Participants, got %d", len(room.Participants))
	}

	// ExpiresAtが24時間後に設定されているかチェック
	expectedExpiry := room.CreatedAt.Add(24 * time.Hour)
	if !room.ExpiresAt.Equal(expectedExpiry) {
		t.Errorf("Expected ExpiresAt to be 24 hours after CreatedAt")
	}
}

func TestRoom_AddParticipant(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)
	userID := uuid.New()

	// 正常な参加者追加
	err := room.AddParticipant(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(room.Participants) != 1 {
		t.Errorf("Expected 1 participant, got %d", len(room.Participants))
	}

	participant := room.Participants[0]
	if participant.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, participant.UserID)
	}
	if participant.IsHost != false {
		t.Errorf("Expected IsHost false, got %v", participant.IsHost)
	}
	if participant.IsMuted != false {
		t.Errorf("Expected IsMuted false, got %v", participant.IsMuted)
	}

	// ホストを追加（IsHostがtrueになることを確認）
	err = room.AddParticipant(hostID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(room.Participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(room.Participants))
	}

	// ホストの参加者を見つける
	var hostParticipant *Participant
	for _, p := range room.Participants {
		if p.UserID == hostID {
			hostParticipant = &p
			break
		}
	}
	if hostParticipant == nil {
		t.Error("Host participant not found")
	} else if !hostParticipant.IsHost {
		t.Error("Expected host participant to have IsHost true")
	}
}

func TestRoom_AddParticipant_DuplicateUser(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)
	userID := uuid.New()

	// 最初の追加は成功
	err := room.AddParticipant(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 同じユーザーの再追加はエラー
	err = room.AddParticipant(userID)
	if err == nil {
		t.Error("Expected error for duplicate user, got nil")
	}
	if err.Error() != "user already in room" {
		t.Errorf("Expected 'user already in room' error, got %v", err)
	}
}

func TestRoom_AddParticipant_RoomFull(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)

	// 定員まで参加者を追加
	for i := 0; i < room.MaxCapacity; i++ {
		userID := uuid.New()
		err := room.AddParticipant(userID)
		if err != nil {
			t.Errorf("Expected no error for participant %d, got %v", i, err)
		}
	}

	// 定員オーバーはエラー
	extraUserID := uuid.New()
	err := room.AddParticipant(extraUserID)
	if err == nil {
		t.Error("Expected error for room full, got nil")
	}
	if err.Error() != "room is full" {
		t.Errorf("Expected 'room is full' error, got %v", err)
	}
}

func TestRoom_AddParticipant_ExpiredRoom(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)
	
	// ルームを期限切れにする
	room.ExpiresAt = time.Now().Add(-1 * time.Hour)
	
	userID := uuid.New()
	err := room.AddParticipant(userID)
	if err == nil {
		t.Error("Expected error for expired room, got nil")
	}
	if err.Error() != "room has expired" {
		t.Errorf("Expected 'room has expired' error, got %v", err)
	}
}

func TestRoom_RemoveParticipant(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)
	userID := uuid.New()

	// 参加者を追加
	room.AddParticipant(userID)
	if len(room.Participants) != 1 {
		t.Errorf("Expected 1 participant, got %d", len(room.Participants))
	}

	// 参加者を削除
	err := room.RemoveParticipant(userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(room.Participants) != 0 {
		t.Errorf("Expected 0 participants, got %d", len(room.Participants))
	}

	// 存在しない参加者の削除はエラー
	err = room.RemoveParticipant(userID)
	if err == nil {
		t.Error("Expected error for non-existent participant, got nil")
	}
	if err.Error() != "participant not found" {
		t.Errorf("Expected 'participant not found' error, got %v", err)
	}
}

func TestRoom_MuteParticipant(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)
	userID := uuid.New()

	// 参加者を追加
	room.AddParticipant(userID)

	// ホストが参加者をミュート
	err := room.MuteParticipant(hostID, userID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	participant, _ := room.GetParticipant(userID)
	if !participant.IsMuted {
		t.Error("Expected participant to be muted")
	}

	// 非ホストがミュートしようとするとエラー
	nonHostID := uuid.New()
	err = room.MuteParticipant(nonHostID, userID)
	if err == nil {
		t.Error("Expected error for non-host muting, got nil")
	}
	if err.Error() != "only host can mute participants" {
		t.Errorf("Expected 'only host can mute participants' error, got %v", err)
	}
}

func TestRoom_IsHost(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)
	userID := uuid.New()

	if !room.IsHost(hostID) {
		t.Error("Expected hostID to be recognized as host")
	}
	if room.IsHost(userID) {
		t.Error("Expected userID to not be recognized as host")
	}
}

func TestRoom_IsParticipant(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)
	userID := uuid.New()

	// 参加前はfalse
	if room.IsParticipant(userID) {
		t.Error("Expected userID to not be participant before joining")
	}

	// 参加後はtrue
	room.AddParticipant(userID)
	if !room.IsParticipant(userID) {
		t.Error("Expected userID to be participant after joining")
	}
}

func TestRoom_IsExpired(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)

	// 新しいルームは期限切れではない
	if room.IsExpired() {
		t.Error("Expected new room to not be expired")
	}

	// 期限を過去に設定
	room.ExpiresAt = time.Now().Add(-1 * time.Hour)
	if !room.IsExpired() {
		t.Error("Expected room with past expiry to be expired")
	}
}

func TestRoom_IsFull(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)

	// 空のルームは満室ではない
	if room.IsFull() {
		t.Error("Expected empty room to not be full")
	}

	// 定員まで追加
	for i := 0; i < room.MaxCapacity; i++ {
		userID := uuid.New()
		room.AddParticipant(userID)
	}

	if !room.IsFull() {
		t.Error("Expected room at capacity to be full")
	}
}

func TestRoom_ExtendExpiry(t *testing.T) {
	hostID := uuid.New()
	room := NewRoom("Test Room", hostID, false)
	originalExpiry := room.ExpiresAt

	// 1時間延長
	extension := 1 * time.Hour
	room.ExtendExpiry(extension)

	expectedExpiry := originalExpiry.Add(extension)
	if !room.ExpiresAt.Equal(expectedExpiry) {
		t.Errorf("Expected ExpiresAt to be %v, got %v", expectedExpiry, room.ExpiresAt)
	}
}
