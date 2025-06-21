package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	googleID := "google123"
	email := "test@example.com"
	name := "Test User"
	avatarURL := "https://example.com/avatar.jpg"

	user := NewUser(googleID, email, name, avatarURL)

	if user.GoogleID != googleID {
		t.Errorf("Expected GoogleID %s, got %s", googleID, user.GoogleID)
	}
	if user.Email != email {
		t.Errorf("Expected Email %s, got %s", email, user.Email)
	}
	if user.Name != name {
		t.Errorf("Expected Name %s, got %s", name, user.Name)
	}
	if user.AvatarURL != avatarURL {
		t.Errorf("Expected AvatarURL %s, got %s", avatarURL, user.AvatarURL)
	}
	if user.ID == uuid.Nil {
		t.Error("Expected ID to be generated")
	}
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		expected bool
	}{
		{
			name: "Valid user",
			user: &User{
				Email: "test@example.com",
				Name:  "Test User",
			},
			expected: true,
		},
		{
			name: "Missing email",
			user: &User{
				Name: "Test User",
			},
			expected: false,
		},
		{
			name: "Missing name",
			user: &User{
				Email: "test@example.com",
			},
			expected: false,
		},
		{
			name: "Empty email and name",
			user: &User{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.IsValid()
			if result != tt.expected {
				t.Errorf("Expected IsValid() to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	user := NewUser("google123", "test@example.com", "Original Name", "original-avatar.jpg")
	originalAvatar := user.AvatarURL

	// 名前のみ更新
	newName := "Updated Name"
	user.UpdateProfile(newName, "")
	
	if user.Name != newName {
		t.Errorf("Expected Name to be updated to %s, got %s", newName, user.Name)
	}
	if user.AvatarURL != originalAvatar {
		t.Errorf("Expected AvatarURL to remain %s, got %s", originalAvatar, user.AvatarURL)
	}

	// アバターのみ更新
	newAvatar := "new-avatar.jpg"
	user.UpdateProfile("", newAvatar)
	
	if user.Name != newName {
		t.Errorf("Expected Name to remain %s, got %s", newName, user.Name)
	}
	if user.AvatarURL != newAvatar {
		t.Errorf("Expected AvatarURL to be updated to %s, got %s", newAvatar, user.AvatarURL)
	}

	// 両方更新
	newerName := "Newer Name"
	newerAvatar := "newer-avatar.jpg"
	user.UpdateProfile(newerName, newerAvatar)
	
	if user.Name != newerName {
		t.Errorf("Expected Name to be updated to %s, got %s", newerName, user.Name)
	}
	if user.AvatarURL != newerAvatar {
		t.Errorf("Expected AvatarURL to be updated to %s, got %s", newerAvatar, user.AvatarURL)
	}

	// 空文字で更新（変更されないことを確認）
	user.UpdateProfile("", "")
	
	if user.Name != newerName {
		t.Errorf("Expected Name to remain %s, got %s", newerName, user.Name)
	}
	if user.AvatarURL != newerAvatar {
		t.Errorf("Expected AvatarURL to remain %s, got %s", newerAvatar, user.AvatarURL)
	}
}

func TestUser_CreatedAtIsRecent(t *testing.T) {
	before := time.Now()
	user := NewUser("google123", "test@example.com", "Test User", "avatar.jpg")
	after := time.Now()

	if user.CreatedAt.Before(before) || user.CreatedAt.After(after) {
		t.Errorf("Expected CreatedAt to be between %v and %v, got %v", before, after, user.CreatedAt)
	}
}
