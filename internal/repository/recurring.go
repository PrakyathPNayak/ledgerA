package repository

import (
	"context"
	"fmt"
	"ledgerA/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type recurringRepository struct {
	db *gorm.DB
}

// NewRecurringRepository creates a new RecurringRepository.
func NewRecurringRepository(db *gorm.DB) RecurringRepository {
	return &recurringRepository{db: db}
}

func (r *recurringRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.RecurringTransaction, int64, error) {
	var items []model.RecurringTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RecurringTransaction{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("recurringRepo.ListByUserID.Count: %w", err)
	}

	if err := query.Order("next_due_date asc").Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("recurringRepo.ListByUserID.Find: %w", err)
	}

	return items, total, nil
}

func (r *recurringRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.RecurringTransaction, error) {
	var item model.RecurringTransaction
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		return nil, fmt.Errorf("recurringRepo.FindByID: %w", err)
	}
	return &item, nil
}

func (r *recurringRepository) Create(ctx context.Context, item *model.RecurringTransaction) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("recurringRepo.Create: %w", err)
	}
	return nil
}

func (r *recurringRepository) Update(ctx context.Context, item *model.RecurringTransaction) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("recurringRepo.Update: %w", err)
	}
	return nil
}

func (r *recurringRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.RecurringTransaction{})
	if res.Error != nil {
		return fmt.Errorf("recurringRepo.Delete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *recurringRepository) FindDue(ctx context.Context, before time.Time) ([]model.RecurringTransaction, error) {
	var items []model.RecurringTransaction
	if err := r.db.WithContext(ctx).
		Where("is_active = ? AND next_due_date <= ? AND (end_date IS NULL OR end_date >= ?)", true, before, before).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("recurringRepo.FindDue: %w", err)
	}
	return items, nil
}
