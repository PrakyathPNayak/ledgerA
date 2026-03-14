package repository

import (
	"context"
	"fmt"
	"ledgerA/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new AccountRepository.
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Account, int64, error) {
	var accounts []model.Account
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Account{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("accountRepo.ListByUserID.Count: %w", err)
	}
	if err := query.Order("created_at desc").Find(&accounts).Error; err != nil {
		return nil, 0, fmt.Errorf("accountRepo.ListByUserID.Find: %w", err)
	}
	return accounts, total, nil
}

func (r *accountRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Account, error) {
	var account model.Account
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		return nil, fmt.Errorf("accountRepo.FindByID: %w", err)
	}
	return &account, nil
}

func (r *accountRepository) Create(ctx context.Context, account *model.Account) error {
	if err := r.db.WithContext(ctx).Create(account).Error; err != nil {
		return fmt.Errorf("accountRepo.Create: %w", err)
	}
	return nil
}

func (r *accountRepository) Update(ctx context.Context, account *model.Account) error {
	if err := r.db.WithContext(ctx).Save(account).Error; err != nil {
		return fmt.Errorf("accountRepo.Update: %w", err)
	}
	return nil
}

func (r *accountRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Account{}).Error; err != nil {
		return fmt.Errorf("accountRepo.Delete: %w", err)
	}
	return nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, userID uuid.UUID, id uuid.UUID, delta float64) error {
	if err := r.db.WithContext(ctx).
		Model(&model.Account{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("current_balance", gorm.Expr("current_balance + ?", delta)).Error; err != nil {
		return fmt.Errorf("accountRepo.UpdateBalance: %w", err)
	}
	return nil
}
