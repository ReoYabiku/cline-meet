package repository

import (
	"context"

	"github.com/cline-meet/backend/internal/domain/model"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type User interface {
	// Create creates a new user
	Create(ctx context.Context, user *model.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	
	// GetByGoogleID retrieves a user by Google ID
	GetByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *model.User) error
	
	// Delete deletes a user
	Delete(ctx context.Context, id uuid.UUID) error
}
