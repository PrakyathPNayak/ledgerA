package repository

import (
	"context"
	"fmt"
	"ledgerA/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new TransactionRepository.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) ListByUserID(ctx context.Context, userID uuid.UUID, filter TransactionListFilter) ([]model.Transaction, int64, error) {
	var txs []model.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Transaction{}).Where("user_id = ?", userID)

	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.SubcategoryID != nil {
		query = query.Where("subcategory_id = ?", *filter.SubcategoryID)
	}
	if filter.DateFrom != nil {
		query = query.Where("transaction_date >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("transaction_date <= ?", *filter.DateTo)
	}
	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("name ILIKE ? OR notes ILIKE ?", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
	}
	if filter.Type == "income" {
		query = query.Where("amount > 0")
	}
	if filter.Type == "expense" {
		query = query.Where("amount < 0")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("transactionRepo.ListByUserID.Count: %w", err)
	}

	if filter.PassbookMode {
		if err := query.Order("transaction_date asc, created_at asc").Find(&txs).Error; err != nil {
			return nil, 0, fmt.Errorf("transactionRepo.ListByUserID.PassbookFind: %w", err)
		}
		return txs, total, nil
	}

	sortBy := "transaction_date"
	sortDir := "desc"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	if filter.SortDir != "" {
		sortDir = filter.SortDir
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

	if err := query.Order(sortBy + " " + sortDir).Offset(offset).Limit(perPage).Find(&txs).Error; err != nil {
		return nil, 0, fmt.Errorf("transactionRepo.ListByUserID.Find: %w", err)
	}
	return txs, total, nil
}

func (r *transactionRepository) FindByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Transaction, error) {
	var tx model.Transaction
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&tx).Error; err != nil {
		return nil, fmt.Errorf("transactionRepo.FindByID: %w", err)
	}
	return &tx, nil
}

func (r *transactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		return fmt.Errorf("transactionRepo.Create: %w", err)
	}
	return nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *model.Transaction) error {
	if err := r.db.WithContext(ctx).Save(tx).Error; err != nil {
		return fmt.Errorf("transactionRepo.Update: %w", err)
	}
	return nil
}

func (r *transactionRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Transaction{}).Error; err != nil {
		return fmt.Errorf("transactionRepo.Delete: %w", err)
	}
	return nil
}
