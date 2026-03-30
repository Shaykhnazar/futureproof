package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents an application user
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	AvatarURL    *string   `json:"avatar_url"`
	AuthProvider string    `json:"auth_provider"` // "email", "google", "github"
	CreatedAt    time.Time `json:"created_at"`
}

// UserProfile contains additional user information
type UserProfile struct {
	UserID       uuid.UUID `json:"user_id"`
	CurrentJobID *uuid.UUID `json:"current_job_id"`
	CityID       *uuid.UUID `json:"city_id"`
	YearsExp     int        `json:"years_exp"`
	Education    *string    `json:"education"`
	TargetJobID  *uuid.UUID `json:"target_job_id"`
	Skills       []string   `json:"skills"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// UserWithProfile combines user and profile data
type UserWithProfile struct {
	User
	Profile *UserProfile `json:"profile,omitempty"`
}

// CreateUserRequest represents user registration payload
type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse contains JWT tokens
type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
