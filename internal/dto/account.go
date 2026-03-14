package dto

import (
	"ledgerA/internal/model"

	"github.com/google/uuid"
)

// CreateAccountRequest contains account creation payload.
type CreateAccountRequest struct {
	Name           string  `json:"name" binding:"required" validate:"required,min=1,max=120"`
	OpeningBalance float64 `json:"opening_balance"`
}

// UpdateAccountRequest contains account update payload.
type UpdateAccountRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=1,max=120"`
}

// AccountResponse contains account data sent to clients.
type AccountResponse struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	AccountType    string    `json:"account_type"`
	OpeningBalance float64   `json:"opening_balance"`
	CurrentBalance float64   `json:"current_balance"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}

// AccountWithBalanceResponse extends account details with last transaction date.
type AccountWithBalanceResponse struct {
	AccountResponse
	LastTransactionDate *string `json:"last_transaction_date,omitempty"`
}

// ToModel converts CreateAccountRequest into model.Account.
func (req CreateAccountRequest) ToModel(userID uuid.UUID) model.Account {
	return model.Account{
		UserID:         userID,
		Name:           req.Name,
		AccountType:    "general",
		OpeningBalance: req.OpeningBalance,
		CurrentBalance: req.OpeningBalance,
	}
}

// NewAccountResponse creates AccountResponse from model.Account.
func NewAccountResponse(account model.Account) AccountResponse {
	return AccountResponse{
		ID:             account.ID,
		UserID:         account.UserID,
		Name:           account.Name,
		AccountType:    account.AccountType,
		OpeningBalance: account.OpeningBalance,
		CurrentBalance: account.CurrentBalance,
		CreatedAt:      account.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      account.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
