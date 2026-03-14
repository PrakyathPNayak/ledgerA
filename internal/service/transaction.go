package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/model"
	"ledgerA/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionService struct {
	db          *gorm.DB
	txRepo      repository.TransactionRepository
	accountRepo repository.AccountRepository
	catRepo     repository.CategoryRepository
	subRepo     repository.SubcategoryRepository
	auditRepo   repository.AuditRepository
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(db *gorm.DB, txRepo repository.TransactionRepository, accountRepo repository.AccountRepository, catRepo repository.CategoryRepository, subRepo repository.SubcategoryRepository, auditRepo repository.AuditRepository) TransactionService {
	return &transactionService{db: db, txRepo: txRepo, accountRepo: accountRepo, catRepo: catRepo, subRepo: subRepo, auditRepo: auditRepo}
}

func (s *transactionService) List(ctx context.Context, userID uuid.UUID, filters dto.TransactionFilters) ([]model.Transaction, int64, error) {
	repoFilter, err := toRepoFilter(filters)
	if err != nil {
		return nil, 0, fmt.Errorf("transactionService.List.Filter: %w", err)
	}
	items, total, err := s.txRepo.ListByUserID(ctx, userID, repoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("transactionService.List: %w", err)
	}
	return items, total, nil
}

func (s *transactionService) Get(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Transaction, error) {
	item, err := s.txRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("transactionService.Get: %w", err)
	}
	return item, nil
}

func (s *transactionService) Create(ctx context.Context, userID uuid.UUID, req dto.CreateTransactionRequest) (*model.Transaction, error) {
	if _, err := s.accountRepo.FindByID(ctx, userID, req.AccountID); err != nil {
		return nil, fmt.Errorf("transactionService.Create.Account: %w", err)
	}

	category, err := s.catRepo.FindByID(ctx, userID, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("transactionService.Create.Category: %w", err)
	}

	var subID uuid.UUID
	if req.SubcategoryID != nil {
		subID = *req.SubcategoryID
	} else {
		sub, subErr := s.subRepo.FindByNameForCategory(ctx, userID, req.CategoryID, category.Name)
		if subErr != nil {
			sub = &model.Subcategory{UserID: userID, CategoryID: req.CategoryID, Name: category.Name}
			if createErr := s.subRepo.Create(ctx, sub); createErr != nil {
				if isDuplicateKeyError(createErr) {
					existing, findErr := s.subRepo.FindByNameForCategory(ctx, userID, req.CategoryID, category.Name)
					if findErr != nil {
						return nil, fmt.Errorf("transactionService.Create.SubcategoryAutoCreate.FindExisting: %w", findErr)
					}
					sub = existing
				} else {
					return nil, fmt.Errorf("transactionService.Create.SubcategoryAutoCreate: %w", createErr)
				}
			}
		}
		subID = sub.ID
	}

	item, err := req.ToModel(userID, subID)
	if err != nil {
		return nil, fmt.Errorf("transactionService.Create.ToModel: %w", err)
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewTransactionRepository(tx)
		accountRepo := repository.NewAccountRepository(tx)
		auditRepo := repository.NewAuditRepository(tx)

		if err := txRepo.Create(ctx, &item); err != nil {
			return fmt.Errorf("transactionService.Create.TxCreate: %w", err)
		}
		if err := accountRepo.UpdateBalance(ctx, userID, req.AccountID, item.Amount); err != nil {
			return fmt.Errorf("transactionService.Create.UpdateBalance: %w", err)
		}
		diffBytes, marshalErr := json.Marshal(map[string]any{"before": nil, "after": item})
		if marshalErr != nil {
			return fmt.Errorf("transactionService.Create.Marshal: %w", marshalErr)
		}
		entry := model.AuditLog{UserID: userID, EntityType: "transaction", EntityID: item.ID, Action: "create", Diff: string(diffBytes)}
		if err := auditRepo.Create(ctx, &entry); err != nil {
			return fmt.Errorf("transactionService.Create.Audit: %w", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("transactionService.Create.Tx: %w", err)
	}

	return &item, nil
}

func (s *transactionService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.UpdateTransactionRequest) (*model.Transaction, error) {
	item, err := s.txRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("transactionService.Update.Find: %w", err)
	}
	before := *item

	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Notes != nil {
		item.Notes = req.Notes
	}
	if req.IsScheduled != nil {
		item.IsScheduled = *req.IsScheduled
	}
	if req.CategoryID != nil {
		item.CategoryID = *req.CategoryID
	}
	if req.SubcategoryID != nil {
		item.SubcategoryID = *req.SubcategoryID
	}
	if req.AccountID != nil {
		item.AccountID = *req.AccountID
	}
	if req.TransactionDate != nil {
		t, parseErr := dto.ParseTransactionDate(*req.TransactionDate)
		if parseErr != nil {
			return nil, fmt.Errorf("transactionService.Update.Date: %w", parseErr)
		}
		item.TransactionDate = t
	}
	if req.Amount != nil {
		item.Amount = *req.Amount
	}

	oldAmount := before.Amount
	newAmount := item.Amount
	delta := newAmount - oldAmount

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewTransactionRepository(tx)
		accountRepo := repository.NewAccountRepository(tx)
		auditRepo := repository.NewAuditRepository(tx)

		if err := txRepo.Update(ctx, item); err != nil {
			return fmt.Errorf("transactionService.Update.TxUpdate: %w", err)
		}
		if delta != 0 {
			if err := accountRepo.UpdateBalance(ctx, userID, item.AccountID, delta); err != nil {
				return fmt.Errorf("transactionService.Update.UpdateBalance: %w", err)
			}
		}
		diffBytes, marshalErr := json.Marshal(map[string]any{"before": before, "after": item})
		if marshalErr != nil {
			return fmt.Errorf("transactionService.Update.Marshal: %w", marshalErr)
		}
		entry := model.AuditLog{UserID: userID, EntityType: "transaction", EntityID: item.ID, Action: "update", Diff: string(diffBytes)}
		if err := auditRepo.Create(ctx, &entry); err != nil {
			return fmt.Errorf("transactionService.Update.Audit: %w", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("transactionService.Update.Tx: %w", err)
	}

	return item, nil
}

func (s *transactionService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	item, err := s.txRepo.FindByID(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("transactionService.Delete.Find: %w", err)
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewTransactionRepository(tx)
		accountRepo := repository.NewAccountRepository(tx)
		auditRepo := repository.NewAuditRepository(tx)

		if err := txRepo.Delete(ctx, userID, id); err != nil {
			return fmt.Errorf("transactionService.Delete.TxDelete: %w", err)
		}
		if err := accountRepo.UpdateBalance(ctx, userID, item.AccountID, -item.Amount); err != nil {
			return fmt.Errorf("transactionService.Delete.UpdateBalance: %w", err)
		}
		diffBytes, marshalErr := json.Marshal(map[string]any{"before": item, "after": nil})
		if marshalErr != nil {
			return fmt.Errorf("transactionService.Delete.Marshal: %w", marshalErr)
		}
		entry := model.AuditLog{UserID: userID, EntityType: "transaction", EntityID: item.ID, Action: "delete", Diff: string(diffBytes)}
		if err := auditRepo.Create(ctx, &entry); err != nil {
			return fmt.Errorf("transactionService.Delete.Audit: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transactionService.Delete.Tx: %w", err)
	}

	return nil
}

func toRepoFilter(filters dto.TransactionFilters) (repository.TransactionListFilter, error) {
	repoFilter := repository.TransactionListFilter{
		AccountID:     filters.AccountID,
		CategoryID:    filters.CategoryID,
		SubcategoryID: filters.SubcategoryID,
		Search:        filters.Search,
		Type:          filters.Type,
		SortBy:        filters.SortBy,
		SortDir:       filters.SortDir,
		Page:          filters.Page,
		PerPage:       filters.PerPage,
		PassbookMode:  filters.PassbookMode,
	}
	if filters.DateFrom != nil && *filters.DateFrom != "" {
		t, err := time.Parse("2006-01-02", *filters.DateFrom)
		if err != nil {
			return repository.TransactionListFilter{}, err
		}
		repoFilter.DateFrom = &t
	}
	if filters.DateTo != nil && *filters.DateTo != "" {
		t, err := time.Parse("2006-01-02", *filters.DateTo)
		if err != nil {
			return repository.TransactionListFilter{}, err
		}
		repoFilter.DateTo = &t
	}
	return repoFilter, nil
}
