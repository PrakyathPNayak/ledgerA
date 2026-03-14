package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog records immutable append-only logs of operations.
type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	EntityType string    `gorm:"type:text;not null;index:idx_audit_entity"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null;index:idx_audit_entity"`
	Action     string    `gorm:"type:text;not null"`
	Diff       string    `gorm:"type:jsonb;not null"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

// BeforeCreate sets a UUID before inserting into the database if not already set.
func (l *AuditLog) BeforeCreate(_ *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
