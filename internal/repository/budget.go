package repository

import (
	"context"
	"fmt"
	"ledgerA/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type budgetRepository struct {
	db *gorm.DB
}

// NewBudgetRepository creates a new BudgetRepository.
func NewBudgetRepository(db *gorm.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Budget, int64, error) {
	var items []model.Budget
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Budget{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("budgetRepo.ListByUserID.Count: %w", err)
	}

	if err := query.Order("created_at desc").Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("budgetRepo.ListByUserID.Find: %w", err)
	}

	return items, total, nil
}

func (r *budgetRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Budget, error) {
	var item model.Budget
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		return nil, fmt.Errorf("budgetRepo.FindByID: %w", err)
	}
	return &item, nil
}

func (r *budgetRepository) Create(ctx context.Context, item *model.Budget) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("budgetRepo.Create: %w", err)
	}
	return nil
}

func (r *budgetRepository) Update(ctx context.Context, item *model.Budget) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("budgetRepo.Update: %w", err)
	}
	return nil
}

func (r *budgetRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Budget{})
	if res.Error != nil {
		return fmt.Errorf("budgetRepo.Delete: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *budgetRepository) FindByCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, period string) (*model.Budget, error) {
	var item model.Budget
	if err := r.db.WithContext(ctx).Where("user_id = ? AND category_id = ? AND period = ?", userID, categoryID, period).First(&item).Error; err != nil {
		return nil, fmt.Errorf("budgetRepo.FindByCategory: %w", err)
	}
	return &item, nil
}
