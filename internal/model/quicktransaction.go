package model

import "github.com/google/uuid"

// QuickTransaction is a template for rapidly creating common transactions.
type QuickTransaction struct {
	BaseModel
	UserID        uuid.UUID  `gorm:"type:uuid;not null"`
	Label         string     `gorm:"type:text;not null"`
	AccountID     *uuid.UUID `gorm:"type:uuid"`
	CategoryID    *uuid.UUID `gorm:"type:uuid"`
	SubcategoryID *uuid.UUID `gorm:"type:uuid"`
	Amount        *float64   `gorm:"type:numeric(20,4)"`
	Notes         *string    `gorm:"type:text"`
	SortOrder     int        `gorm:"type:integer;not null;default:0"`
}
