package usecase

import (
	"context"
	"fmt"

	"github.com/cline-meet/backend/internal/domain/model"
	"github.com/cline-meet/backend/internal/domain/repository"
	"github.com/cline-meet/backend/internal/domain/service"
	"github.com/google/uuid"
)

// User handles user-related business logic
type User struct {
	userRepo         repository.User
	realtimeNotifier service.RealtimeNotifier
	sessionManager   service.SessionManager
}

// NewUser creates a new User usecase
func NewUser(
	userRepo repository.User,
	realtimeNotifier service.RealtimeNotifier,
	sessionManager service.SessionManager,
) *User {
	return &User{
		userRepo:         userRepo,
		realtimeNotifier: realtimeNotifier,
		sessionManager:   sessionManager,
	}
}

// CreateUser creates a new user
func (u *User) CreateUser(ctx context.Context, googleID, email, name, avatarURL string) (*model.User, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByGoogleID(ctx, googleID)
	if err == nil {
		// User already exists, return existing user
		return existingUser, nil
	}

	// Create new user
	user := model.NewUser(googleID, email, name, avatarURL)

	// Validate user data
	if !user.IsValid() {
		return nil, fmt.Errorf("invalid user data: email and name are required")
	}

	// Save to repository
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (u *User) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

// GetUserByGoogleID retrieves a user by Google ID
func (u *User) GetUserByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	user, err := u.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (u *User) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(ctx context.Context, userID uuid.UUID, name, avatarURL string) error {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update profile
	user.UpdateProfile(name, avatarURL)

	// Validate updated data
	if !user.IsValid() {
		return fmt.Errorf("invalid user data: email and name are required")
	}

	// Save to repository
	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user
func (u *User) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Check if user exists
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Delete user session if exists
	if err := u.sessionManager.DeleteSession(ctx, userID); err != nil {
		// Log error but don't fail the delete operation
		// Session deletion failure shouldn't prevent user deletion
	}

	// Delete user from repository
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// LoginUser handles user login process
func (u *User) LoginUser(ctx context.Context, googleID, email, name, avatarURL string) (*model.User, error) {
	// Try to get existing user
	user, err := u.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		// User doesn't exist, create new user
		return u.CreateUser(ctx, googleID, email, name, avatarURL)
	}

	// User exists, update profile if needed
	if user.Name != name || user.AvatarURL != avatarURL {
		user.UpdateProfile(name, avatarURL)
		if err := u.userRepo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update user profile: %w", err)
		}
	}

	return user, nil
}
