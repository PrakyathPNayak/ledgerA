package dto

import (
	"ledgerA/internal/model"

	"github.com/google/uuid"
)

// CreateBudgetRequest contains budget creation payload.
type CreateBudgetRequest struct {
	CategoryID uuid.UUID `json:"category_id" binding:"required" validate:"required"`
	Amount     float64   `json:"amount" binding:"required" validate:"required,gt=0"`
	Period     string    `json:"period" binding:"required" validate:"required,oneof=monthly yearly"`
}

// UpdateBudgetRequest contains budget update payload.
type UpdateBudgetRequest struct {
	Amount   *float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// BudgetResponse contains budget data sent to clients.
type BudgetResponse struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	Amount     float64   `json:"amount"`
	Period     string    `json:"period"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

// BudgetProgressResponse extends budget with spent amount and remaining.
type BudgetProgressResponse struct {
	BudgetResponse
	Spent     float64 `json:"spent"`
	Remaining float64 `json:"remaining"`
	Percent   float64 `json:"percent"`
}

// NewBudgetResponse creates BudgetResponse from model.
func NewBudgetResponse(b model.Budget) BudgetResponse {
	return BudgetResponse{
		ID:         b.ID,
		UserID:     b.UserID,
		CategoryID: b.CategoryID,
		Amount:     b.Amount,
		Period:     b.Period,
		IsActive:   b.IsActive,
		CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
