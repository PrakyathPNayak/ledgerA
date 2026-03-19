package model

import (
	"time"

	"github.com/google/uuid"
)

// RecurringTransaction represents a scheduled recurring transaction template.
type RecurringTransaction struct {
	BaseModel
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index:idx_recurring_user,where:deleted_at IS NULL"`
	AccountID      uuid.UUID  `gorm:"type:uuid;not null"`
	CategoryID     uuid.UUID  `gorm:"type:uuid;not null"`
	SubcategoryID  uuid.UUID  `gorm:"type:uuid;not null"`
	Name           string     `gorm:"type:text;not null"`
	Amount         float64    `gorm:"type:numeric(20,4);not null"`
	Notes          *string    `gorm:"type:text"`
	Frequency      string     `gorm:"type:text;not null"` // daily, weekly, monthly, yearly
	StartDate      time.Time  `gorm:"type:date;not null"`
	NextDueDate    time.Time  `gorm:"type:date;not null;index:idx_recurring_due,where:deleted_at IS NULL AND is_active = true"`
	EndDate        *time.Time `gorm:"type:date"`
	IsActive       bool       `gorm:"type:boolean;not null;default:true"`
	LastExecutedAt *time.Time `gorm:"type:timestamptz"`
}
