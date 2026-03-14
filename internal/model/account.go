package model

import "github.com/google/uuid"

// Account represents a financial account (e.g., bank, wallet) for a user.
type Account struct {
	BaseModel
	UserID         uuid.UUID `gorm:"type:uuid;not null;index:idx_accounts_user_id,where:deleted_at IS NULL"`
	Name           string    `gorm:"type:text;not null"`
	AccountType    string    `gorm:"type:text;not null;default:'general'"`
	OpeningBalance float64   `gorm:"type:numeric(20,4);not null;default:0"`
	CurrentBalance float64   `gorm:"type:numeric(20,4);not null;default:0"`
}
