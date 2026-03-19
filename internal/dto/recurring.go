package dto

import (
	"ledgerA/internal/model"
	"time"

	"github.com/google/uuid"
)

// CreateRecurringRequest contains recurring transaction creation payload.
type CreateRecurringRequest struct {
	AccountID     uuid.UUID  `json:"account_id" binding:"required" validate:"required"`
	CategoryID    uuid.UUID  `json:"category_id" binding:"required" validate:"required"`
	SubcategoryID *uuid.UUID `json:"subcategory_id,omitempty"`
	Name          string     `json:"name" binding:"required" validate:"required,min=1,max=200"`
	Amount        float64    `json:"amount" binding:"required" validate:"required"`
	Notes         *string    `json:"notes,omitempty" validate:"omitempty,max=2000"`
	Frequency     string     `json:"frequency" binding:"required" validate:"required,oneof=daily weekly monthly yearly"`
	StartDate     string     `json:"start_date" binding:"required" validate:"required,datetime=2006-01-02"`
	EndDate       *string    `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// UpdateRecurringRequest contains recurring transaction update payload.
type UpdateRecurringRequest struct {
	AccountID     *uuid.UUID `json:"account_id,omitempty"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	SubcategoryID *uuid.UUID `json:"subcategory_id,omitempty"`
	Name          *string    `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Amount        *float64   `json:"amount,omitempty"`
	Notes         *string    `json:"notes,omitempty" validate:"omitempty,max=2000"`
	Frequency     *string    `json:"frequency,omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
	EndDate       *string    `json:"end_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

// RecurringResponse contains recurring transaction data sent to clients.
type RecurringResponse struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	AccountID      uuid.UUID  `json:"account_id"`
	CategoryID     uuid.UUID  `json:"category_id"`
	SubcategoryID  uuid.UUID  `json:"subcategory_id"`
	Name           string     `json:"name"`
	Amount         float64    `json:"amount"`
	Notes          *string    `json:"notes,omitempty"`
	Frequency      string     `json:"frequency"`
	StartDate      string     `json:"start_date"`
	NextDueDate    string     `json:"next_due_date"`
	EndDate        *string    `json:"end_date,omitempty"`
	IsActive       bool       `json:"is_active"`
	LastExecutedAt *string    `json:"last_executed_at,omitempty"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
}

// NewRecurringResponse creates RecurringResponse from model.
func NewRecurringResponse(r model.RecurringTransaction) RecurringResponse {
	resp := RecurringResponse{
		ID:            r.ID,
		UserID:        r.UserID,
		AccountID:     r.AccountID,
		CategoryID:    r.CategoryID,
		SubcategoryID: r.SubcategoryID,
		Name:          r.Name,
		Amount:        r.Amount,
		Notes:         r.Notes,
		Frequency:     r.Frequency,
		StartDate:     r.StartDate.Format("2006-01-02"),
		NextDueDate:   r.NextDueDate.Format("2006-01-02"),
		IsActive:      r.IsActive,
		CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if r.EndDate != nil {
		s := r.EndDate.Format("2006-01-02")
		resp.EndDate = &s
	}
	if r.LastExecutedAt != nil {
		s := r.LastExecutedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.LastExecutedAt = &s
	}
	return resp
}

// ComputeNextDueDate computes next due date after the given date based on frequency.
func ComputeNextDueDate(current time.Time, frequency string) time.Time {
	switch frequency {
	case "daily":
		return current.AddDate(0, 0, 1)
	case "weekly":
		return current.AddDate(0, 0, 7)
	case "monthly":
		return current.AddDate(0, 1, 0)
	case "yearly":
		return current.AddDate(1, 0, 0)
	default:
		return current.AddDate(0, 1, 0)
	}
}
