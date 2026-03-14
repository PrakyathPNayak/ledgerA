package dto

import (
	"fmt"
	"ledgerA/internal/model"
	"time"

	"github.com/google/uuid"
)

// CreateTransactionRequest contains transaction creation payload.
type CreateTransactionRequest struct {
	AccountID       uuid.UUID  `json:"account_id" binding:"required" validate:"required"`
	CategoryID      uuid.UUID  `json:"category_id" binding:"required" validate:"required"`
	SubcategoryID   *uuid.UUID `json:"subcategory_id,omitempty"`
	Name            string     `json:"name" binding:"required" validate:"required,min=1,max=200"`
	Amount          float64    `json:"amount" binding:"required" validate:"required"`
	TransactionDate string     `json:"transaction_date" binding:"required" validate:"required,datetime=2006-01-02"`
	Notes           *string    `json:"notes,omitempty" validate:"omitempty,max=2000"`
	IsScheduled     bool       `json:"is_scheduled"`
}

// UpdateTransactionRequest contains transaction update payload.
type UpdateTransactionRequest struct {
	AccountID       *uuid.UUID `json:"account_id,omitempty"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	SubcategoryID   *uuid.UUID `json:"subcategory_id,omitempty"`
	Name            *string    `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Amount          *float64   `json:"amount,omitempty"`
	TransactionDate *string    `json:"transaction_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Notes           *string    `json:"notes,omitempty" validate:"omitempty,max=2000"`
	IsScheduled     *bool      `json:"is_scheduled,omitempty"`
}

// TransactionResponse contains transaction data sent to clients.
type TransactionResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	AccountID       uuid.UUID `json:"account_id"`
	CategoryID      uuid.UUID `json:"category_id"`
	SubcategoryID   uuid.UUID `json:"subcategory_id"`
	Name            string    `json:"name"`
	Amount          float64   `json:"amount"`
	TransactionDate string    `json:"transaction_date"`
	Notes           *string   `json:"notes,omitempty"`
	IsScheduled     bool      `json:"is_scheduled"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

// TransactionListResponse contains items and total count.
type TransactionListResponse struct {
	Items []TransactionResponse `json:"items"`
	Total int64                 `json:"total"`
}

// TransactionFilters contains list filters.
type TransactionFilters struct {
	AccountID     *uuid.UUID `form:"account_id"`
	CategoryID    *uuid.UUID `form:"category_id"`
	SubcategoryID *uuid.UUID `form:"subcategory_id"`
	DateFrom      *string    `form:"date_from"`
	DateTo        *string    `form:"date_to"`
	Search        *string    `form:"search"`
	Type          string     `form:"type" validate:"omitempty,oneof=income expense all"`
	SortBy        string     `form:"sort_by" validate:"omitempty,oneof=transaction_date amount name created_at"`
	SortDir       string     `form:"sort_dir" validate:"omitempty,oneof=asc desc"`
	Page          int        `form:"page" validate:"omitempty,min=1"`
	PerPage       int        `form:"per_page" validate:"omitempty,min=1,max=200"`
	PassbookMode  bool       `form:"passbook_mode"`
}

// ToModel converts CreateTransactionRequest to model.Transaction.
func (req CreateTransactionRequest) ToModel(userID uuid.UUID, subcategoryID uuid.UUID) (model.Transaction, error) {
	txDate, err := ParseTransactionDate(req.TransactionDate)
	if err != nil {
		return model.Transaction{}, err
	}

	return model.Transaction{
		UserID:          userID,
		AccountID:       req.AccountID,
		CategoryID:      req.CategoryID,
		SubcategoryID:   subcategoryID,
		Name:            req.Name,
		Amount:          req.Amount,
		TransactionDate: txDate,
		Notes:           req.Notes,
		IsScheduled:     req.IsScheduled,
	}, nil
}

// ParseTransactionDate parses accepted transaction date formats.
func ParseTransactionDate(value string) (time.Time, error) {
	layouts := []string{"2006-01-02", "02/01/2006", "01/02/2006"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid transaction date format: %s", value)
}

// NewTransactionResponse creates TransactionResponse from model.Transaction.
func NewTransactionResponse(tx model.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:              tx.ID,
		UserID:          tx.UserID,
		AccountID:       tx.AccountID,
		CategoryID:      tx.CategoryID,
		SubcategoryID:   tx.SubcategoryID,
		Name:            tx.Name,
		Amount:          tx.Amount,
		TransactionDate: tx.TransactionDate.Format("2006-01-02"),
		Notes:           tx.Notes,
		IsScheduled:     tx.IsScheduled,
		CreatedAt:       tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       tx.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
