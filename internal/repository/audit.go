package repository

import (
	"context"
	"fmt"
	"ledgerA/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new AuditRepository.
func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(ctx context.Context, entry *model.AuditLog) error {
	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return fmt.Errorf("auditRepo.Create: %w", err)
	}
	return nil
}

func (r *auditRepository) ListByUserID(ctx context.Context, userID uuid.UUID, filter AuditListFilter) ([]model.AuditLog, int64, error) {
	var items []model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{}).Where("user_id = ?", userID)
	if filter.EntityType != nil {
		query = query.Where("entity_type = ?", *filter.EntityType)
	}
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("auditRepo.ListByUserID.Count: %w", err)
	}

	page := filter.Page
	perPage := filter.PerPage
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	if err := query.Order("created_at desc").Offset(offset).Limit(perPage).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("auditRepo.ListByUserID.Find: %w", err)
	}
	return items, total, nil
}
