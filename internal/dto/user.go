package dto

import (
	"ledgerA/internal/model"

	"github.com/google/uuid"
)

// UpdateUserRequest contains fields that can be updated on user profile.
type UpdateUserRequest struct {
	DisplayName string `json:"display_name" binding:"required" validate:"required,min=1,max=100"`
}

// UserResponse contains user profile data sent to clients.
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	FirebaseUID  string    `json:"firebase_uid"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	CurrencyCode string    `json:"currency_code"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

// NewUserResponse builds UserResponse from a model.User.
func NewUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:           user.ID,
		FirebaseUID:  user.FirebaseUID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		CurrencyCode: user.CurrencyCode,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
