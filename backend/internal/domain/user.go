package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	GoogleID  string    `json:"googleId"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewUser creates a new user
func NewUser(googleID, email, name, avatarURL string) *User {
	return &User{
		ID:        uuid.New(),
		GoogleID:  googleID,
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
		CreatedAt: time.Now(),
	}
}

// IsValid validates user data
func (u *User) IsValid() bool {
	return u.Email != "" && u.Name != ""
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name, avatarURL string) {
	if name != "" {
		u.Name = name
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
}
