package repository

import (
	"context"
	"fmt"
	"ledgerA/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type quickTransactionRepository struct {
	db *gorm.DB
}

// NewQuickTransactionRepository creates a new QuickTransactionRepository.
func NewQuickTransactionRepository(db *gorm.DB) QuickTransactionRepository {
	return &quickTransactionRepository{db: db}
}

func (r *quickTransactionRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.QuickTransaction, int64, error) {
	var quick []model.QuickTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&model.QuickTransaction{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("quickTransactionRepo.ListByUserID.Count: %w", err)
	}
	if err := query.Order("sort_order asc, created_at asc").Find(&quick).Error; err != nil {
		return nil, 0, fmt.Errorf("quickTransactionRepo.ListByUserID.Find: %w", err)
	}
	return quick, total, nil
}

func (r *quickTransactionRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.QuickTransaction, error) {
	var quick model.QuickTransaction
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&quick).Error; err != nil {
		return nil, fmt.Errorf("quickTransactionRepo.FindByID: %w", err)
	}
	return &quick, nil
}

func (r *quickTransactionRepository) Create(ctx context.Context, quick *model.QuickTransaction) error {
	if err := r.db.WithContext(ctx).Create(quick).Error; err != nil {
		return fmt.Errorf("quickTransactionRepo.Create: %w", err)
	}
	return nil
}

func (r *quickTransactionRepository) Update(ctx context.Context, quick *model.QuickTransaction) error {
	if err := r.db.WithContext(ctx).Save(quick).Error; err != nil {
		return fmt.Errorf("quickTransactionRepo.Update: %w", err)
	}
	return nil
}

func (r *quickTransactionRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.QuickTransaction{}).Error; err != nil {
		return fmt.Errorf("quickTransactionRepo.Delete: %w", err)
	}
	return nil
}

func (r *quickTransactionRepository) Reorder(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for idx, id := range ids {
			if err := tx.Model(&model.QuickTransaction{}).
				Where("id = ? AND user_id = ?", id, userID).
				Update("sort_order", idx).Error; err != nil {
				return fmt.Errorf("quickTransactionRepo.Reorder.Update: %w", err)
			}
		}
		return nil
	})
}
