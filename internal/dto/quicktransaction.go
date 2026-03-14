package dto

import (
	"ledgerA/internal/model"

	"github.com/google/uuid"
)

// CreateQuickTransactionRequest contains quick transaction creation payload.
type CreateQuickTransactionRequest struct {
	Label         string     `json:"label" binding:"required" validate:"required,min=1,max=100"`
	AccountID     *uuid.UUID `json:"account_id,omitempty"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	SubcategoryID *uuid.UUID `json:"subcategory_id,omitempty"`
	Amount        *float64   `json:"amount,omitempty"`
	Notes         *string    `json:"notes,omitempty" validate:"omitempty,max=2000"`
}

// UpdateQuickTransactionRequest contains quick transaction update payload.
type UpdateQuickTransactionRequest struct {
	Label         *string    `json:"label,omitempty" validate:"omitempty,min=1,max=100"`
	AccountID     *uuid.UUID `json:"account_id,omitempty"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	SubcategoryID *uuid.UUID `json:"subcategory_id,omitempty"`
	Amount        *float64   `json:"amount,omitempty"`
	Notes         *string    `json:"notes,omitempty" validate:"omitempty,max=2000"`
}

// QuickTransactionResponse contains quick transaction data sent to clients.
type QuickTransactionResponse struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Label         string     `json:"label"`
	AccountID     *uuid.UUID `json:"account_id,omitempty"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	SubcategoryID *uuid.UUID `json:"subcategory_id,omitempty"`
	Amount        *float64   `json:"amount,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	SortOrder     int        `json:"sort_order"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
}

// ReorderQuickTransactionsRequest contains ids in new order.
type ReorderQuickTransactionsRequest struct {
	IDs []uuid.UUID `json:"ids" binding:"required" validate:"required,min=1,dive,required"`
}

// ToModel converts CreateQuickTransactionRequest into model.QuickTransaction.
func (req CreateQuickTransactionRequest) ToModel(userID uuid.UUID) model.QuickTransaction {
	return model.QuickTransaction{
		UserID:        userID,
		Label:         req.Label,
		AccountID:     req.AccountID,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		Amount:        req.Amount,
		Notes:         req.Notes,
	}
}

// NewQuickTransactionResponse creates QuickTransactionResponse from model.QuickTransaction.
func NewQuickTransactionResponse(q model.QuickTransaction) QuickTransactionResponse {
	return QuickTransactionResponse{
		ID:            q.ID,
		UserID:        q.UserID,
		Label:         q.Label,
		AccountID:     q.AccountID,
		CategoryID:    q.CategoryID,
		SubcategoryID: q.SubcategoryID,
		Amount:        q.Amount,
		Notes:         q.Notes,
		SortOrder:     q.SortOrder,
		CreatedAt:     q.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     q.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
