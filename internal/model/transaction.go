package model

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a single income or expense record.
type Transaction struct {
	BaseModel
	UserID          uuid.UUID `gorm:"type:uuid;not null;index:idx_txn_user_date,where:deleted_at IS NULL"`
	AccountID       uuid.UUID `gorm:"type:uuid;not null;index:idx_txn_account,where:deleted_at IS NULL"`
	CategoryID      uuid.UUID `gorm:"type:uuid;not null;index:idx_txn_category,where:deleted_at IS NULL"`
	SubcategoryID   uuid.UUID `gorm:"type:uuid;not null"`
	Name            string    `gorm:"type:text;not null"`
	Amount          float64   `gorm:"type:numeric(20,4);not null"`
	TransactionDate time.Time `gorm:"type:date;not null;index:idx_txn_user_date,where:deleted_at IS NULL"`
	Notes           *string   `gorm:"type:text"`
	IsScheduled     bool      `gorm:"type:boolean;not null;default:false"`
}
